# encoding: utf-8
#
# Encoon : data structuration, presentation and navigation.
# 
# Copyright (C) 2012 David Lambert
# 
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
# 
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
# 
# See doc/COPYRIGHT.rdoc for more details.
class Grid < Entity
  PHASE_CREATE = 'create'
  PHASE_NEW_VERSION = 'new_version'
  PHASE_UPDATE = 'update'

  DISPLAY_ROWS_LIMIT = 10
  DISPLAY_ROWS_LIMIT_FULL = 50

  belongs_to :workspace, :foreign_key => "workspace_uuid", :primary_key => "uuid"
  has_many :grid_locs, :foreign_key => "uuid", :primary_key => "uuid"
  has_many :columns, :foreign_key => "grid_uuid", :primary_key => "uuid"
  has_many :column_locs, :foreign_key => "uuid", :primary_key => "uuid", :through => :columns
  has_many :column_mappings, :foreign_key => "uuid", :primary_key => "uuid", :through => :columns
  has_many :rows, :foreign_key => "grid_uuid", :primary_key => "uuid"
  has_many :grid_mappings, :foreign_key => "grid_uuid", :primary_key => "uuid"
  validates_presence_of :workspace_uuid
  validates_associated :workspace
  attr_reader :loaded, :column_all, :filtered_columns, :has_translation
  
  @loaded = @can_select_data = @can_create_data = @can_update_data = @has_translation = false

  # Loads in memory the structure of the grid
  # and all the information to render the grid: workspace, columns
  # definition, table mapping, column mapping and security settings.
  # Load is mandatory before any access to the data using the grid.
  def load(filters=nil, skip_mapping=false)
    log_debug "Grid#load_grid_structure(#{filters.inspect}, #{skip_mapping}) [#{to_s}]"
    load_workspace
    load_mapping
    load_columns(filters, false, skip_mapping)
    @loaded = true
  end

  # Selects data based on its uuid in the given collection.
  def self.select_entity_by_uuid(collection, uuid)
    log_debug "Grid#select_entity_by_uuid(#{collection}, #{uuid}) [#{to_s}]"
    collection.
      select(self.all_select_columns).
      joins(:grid_locs).
      where("grids.uuid = ?", uuid).
      where("grid_locs.version = grids.version").
      where(as_of_date_clause("grids")).
      where(locale_clause("grid_locs")).
      where(grid_security_clause("grids")).
      first 
  end
  
  # Selects data based on workspace and uri in the given collection.
  def self.select_entity_by_workspace_and_uri(collection, workspace_uuid, uri)
    log_debug "Grid#select_entity_by_workspace_and_uri(#{collection}, #{uri}) [#{to_s}]"
    collection.
      select(self.all_select_columns).
      joins(:grid_locs).
      where("grids.workspace_uuid = ?", workspace_uuid).
      where("grids.uri = ?", uri).
      where("grid_locs.version = grids.version").
      where(as_of_date_clause("grids")).
      where(locale_clause("grid_locs")).
      where(grid_security_clause("grids")).
      first 
  end
  
  # Selects data based on its uuid in the given collection and for a given version number.
  def self.select_entity_by_uuid_version(collection, uuid, version)
    log_debug "Grid#select_entity_by_uuid_version(#{collection}, #{uuid}, #{version}) [#{to_s}]"
    collection.find(:first, 
                    :joins => :grid_locs,
                    :select => self.all_select_columns,
                    :conditions => 
                        ["grids.uuid = :uuid " + 
                         " AND grids.version = :version " + 
                         " AND grid_locs.version = grids.version " +
                         " AND " + locale_clause("grid_locs") + 
                         " AND " + grid_security_clause("grids"), 
                        {:uuid => uuid, :version => version}]) 
  end
  
  # Returns the name of the grid used as a reference
  def reference_name
    name + (workspace.present? ? " [" + workspace_name + "]" : "")
  end

  # Selects the reference rows attached to one grid
  # This is used to display a drop-down list
  def self.select_reference_rows(grid_uuid)
    log_debug "Grid#select_reference_rows(#{grid_uuid}) [#{to_s}]"
    grid = Grid.select_entity_by_uuid(Grid, grid_uuid)
    unless grid.nil?
      grid.load
      grid.row_all(nil, nil, -1)
    end
  end

  # Selects the name of one reference row attached to one grid
  def select_reference_row_name(row_uuid)
    log_debug "Grid#select_reference_row_name(#{row_uuid}) [#{to_s}]"
    row = row_select_entity_by_uuid(row_uuid)
    row.present? ? row_title(row) : ""
  end

  # Selects the description of one reference row attached to one grid
  def select_reference_row_description(row_uuid)
    log_debug "Grid#select_reference_row_description(#{row_uuid}) [#{to_s}]"
    row = row_select_entity_by_uuid(row_uuid)
    row.present? ? row.description : ""
  end

  def workspace_name
    workspace = Workspace.select_entity_by_uuid(Workspace, self.workspace_uuid)
    workspace.present? ? workspace.name : ""
  end

  def select_grid_cast(row_uuid)
    log_debug "Grid#select_grid_cast(#{row_uuid}) [#{to_s}]"
    if self.uuid == GRID_UUID
      grid = Grid::select_entity_by_uuid(Grid, row_uuid)
      grid.load if grid.present?
      grid
    end
  end

  def filter_column_uuid
    attribute_present?(:filter_column_uuid) ? read_attribute(:filter_column_uuid) : nil
  end

  def filter_column_name
    attribute_present?(:filter_column_name) ? read_attribute(:filter_column_name) : nil
  end

  # Builds a filter based on the given row used to select dependent tables.
  def get_filters_on_row(row, row_name)
    self.filter_column_uuid.nil? ? [] : [{:column_uuid => self.filter_column_uuid,
                                          :row_uuid => row.uuid,
                                          :row_name => row_name}]
  end

  # Selects the grids attached to one row via one or more columns.
  def self.select_referenced_grids(uuid)
    log_debug "Grid#select_referenced_grids(uuid=#{uuid}) [#{to_s}]"
    Grid.find_by_sql(["SELECT grids.id," + 
                      " grids.uuid," + 
                      " grids.uri," + 
                      " grid_locs.name," + 
                      " grid_locs.description," + 
                      " grids.version," + 
                      " grids.begin," + 
                      " grids.end," + 
                      " grids.workspace_uuid," + 
                      " grids.enabled," + 
                      " grids.has_name," + 
                      " grids.has_description," + 
                      " grids.create_user_uuid," + 
                      " columns.uuid as filter_column_uuid," +
                      " column_locs.name as filter_column_name" +
                      " FROM grids, grid_locs, columns, column_locs" +
                      " WHERE " + as_of_date_clause("grids") + 
                      " AND grids.uuid = grid_locs.uuid" + 
                      " AND grids.version = grid_locs.version" +
                      " AND " + locale_clause("grid_locs") +
                      " AND columns.grid_uuid = grids.uuid" + 
                      " AND " + as_of_date_clause("columns") + 
                      " AND columns.uuid = column_locs.uuid" + 
                      " AND columns.version = column_locs.version" +
                      " AND " + locale_clause("column_locs") +
                      " AND columns.kind = :kind" + 
                      " AND columns.grid_reference_uuid = :uuid" + 
                      " AND grid_locs.version = grids.version" +
                      " AND " + grid_security_clause("grids") + 
                      " ORDER BY grid_locs.name, column_locs.name", 
                       {:kind => COLUMN_TYPE_REFERENCE, :uuid => uuid}])
  end

  def column_select_by_uuid(uuid)
    self.column_all.each{|column| return column if column.uuid == uuid}
  end

  # Selects data based on uuid
  def column_select_entity_by_uuid(uuid)
    log_debug "Grid#column_select_entity_by_uuid(uuid=#{uuid}) [#{to_s}]"
    columns.find(:first, 
                 :joins => :column_locs,
                 :select => column_all_select_columns,
                 :conditions => 
                      ["columns.uuid = :uuid " + 
                       "AND " + as_of_date_clause("columns") + 
                       "AND column_locs.version = columns.version " +
                       "AND " + locale_clause("column_locs"), 
                      {:uuid => uuid}]) 
  end

  def column_select_entity_by_uuid_version(uuid, version)
    log_debug "Grid#column_select_entity_by_uuid_version(#{uuid}, #{version}) [#{to_s}]"
    columns.find(:first, 
                 :joins => :column_locs,
                 :select => column_all_select_columns,
                 :conditions => 
                      ["columns.uuid = :uuid " + 
                       "AND columns.version = :version " + 
                       "AND column_locs.version = columns.version " +
                       "AND " + locale_clause("column_locs"), 
                      {:uuid => uuid, :version => version}]) 
  end

  # Selects data based on id
  def column_select_entity_by_id(id)
    log_debug "Grid#column_select_entity_by_id(id=#{id}) [#{to_s}]"
    columns.find(:first, 
                 :joins => :column_locs,
                 :select => column_all_select_columns,
                 :conditions => 
                       ["columns.id = :id " +
                        "AND column_locs.version = columns.version " + 
                        "AND " + locale_clause("column_locs"), 
                       {:id => id}]) 
  end

  def column_all_versions(uuid)
    log_debug "Grid#column_all_versions(uuid=#{uuid}) [#{to_s}]"
    columns.find(:all, 
                 :joins => :column_locs,
                 :select => column_all_select_columns,
                 :conditions => 
                        ["columns.uuid = :uuid " +
                         "AND column_locs.version = columns.version " + 
                         "AND " + locale_clause("column_locs"), 
                        {:uuid => uuid}], 
                 :order => "columns.begin")
  end

  def row_count(filters=nil)
    row_all(filters, nil, nil, false, true)
  end

  def row_all(filters=nil, search=nil, page=nil, full=false, count=false)
    log_debug "Grid#row_all(filters=#{filters.inspect}, " +
              "search=#{search}, page=#{page}, " +
              "count=#{count}) [#{to_s}]"
    conditions = ""
    if filters.present?
      filters.each do |filter|
        column_uuid = filter[:column_uuid]
        row_uuid = filter[:row_uuid]
        column_all.each do |column|
          log_debug "Grid#row_all filter: column.uuid=#{column.uuid}"
          if column.uuid == column_uuid and column.kind == COLUMN_TYPE_REFERENCE
            log_debug "Grid#row_all found reference"
            conditions << " AND rows.#{column.physical_column} = #{quote(row_uuid)}"
          end
        end
      end
    end
    if search.present? and @has_translation
      conditions << " AND ("
      conditions << " lower(row_locs.name) like #{quote('%%' + search + '%%')}"
      conditions << " OR lower(row_locs.description) like #{quote('%%' + search + '%%')}"
      conditions << " )"
    end
    offset = 0
    limit = full ? DISPLAY_ROWS_LIMIT_FULL : DISPLAY_ROWS_LIMIT
    if page.present? and page.to_i > 0
      offset = (page.to_i - 1) * limit
    elsif page.present? and page.to_i < 0
      limit = 10000000
    end
    sql = "SELECT " + 
          (count ? "count(*)" : "#{row_all_select_columns}") + 
          " FROM grids grids, #{@db_table} rows" + 
          (@has_translation ? ", #{@db_loc_table} row_locs" : "") +
          " WHERE grids.uuid = #{quote(self.uuid)}" +
          " AND " + as_of_date_clause("grids") +
          " AND " + Grid::grid_security_clause("grids") + 
          " AND " + as_of_date_clause("rows") +
          (@has_mapping ? "" : " AND rows.grid_uuid = #{quote(self.uuid)}") +
          (@has_translation ? 
              " AND rows.uuid = row_locs.uuid" +
              " AND rows.version = row_locs.version" +
              " AND " + locale_clause("row_locs") : "") +
          ((self.uuid == WORKSPACE_UUID) ? 
              " AND " + Grid::workspace_security_clause("rows", false, false) : "") + 
          ((self.uuid == GRID_UUID) ? 
              " AND " + Grid::grid_security_clause("rows") : "") + 
          conditions +
          ((@has_translation and not count) ? " ORDER BY row_locs.name" : "") +
          (count ? "" : " LIMIT #{offset}, #{limit}") 
    count ? connection.select_value(sql).to_i : Row.find_by_sql([sql])
    rescue ActiveRecord::StatementInvalid => exception
      log_error "Grid#row_all #{exception.to_s}" +
                ", data grid='#{name}', sql=#{sql}"
      count ? 0 : []
  end

  # Selects a row based on its uuid or uri.
  def row_select_entity_by_uuid(uuid, uri=nil, version=nil)
    log_debug "Grid#row_select_entity_by_uuid(#{uuid}, #{uri}, #{version}) [#{to_s}]"
    if not @can_select_data
      log_security_warning "Grid#row_select_entity_by_uuid Can't select data"
      return nil
    end
    sql = "SELECT #{row_all_select_columns}" +
          " FROM grids grids, #{@db_table} rows" +
          (@has_translation ? ", #{@db_loc_table} row_locs" : "") +
          " WHERE grids.uuid = #{quote(self.uuid)}" +
          " AND " + as_of_date_clause("grids") +
          " AND " + Grid::grid_security_clause("grids") +
          (uri.nil? ? " AND rows.uuid = #{quote(uuid)}" : " AND rows.uri = #{quote(uri)}") +
          (version.nil? ? " AND " + as_of_date_clause("rows") : " AND rows.version = #{version}") +
          (@has_mapping ? "" : " AND rows.grid_uuid = #{quote(self.uuid)}") +
          (@has_translation ? 
              " AND rows.uuid = row_locs.uuid" +
              " AND rows.version = row_locs.version" +
              " AND " + locale_clause("row_locs") : "") +
          ((self.uuid == WORKSPACE_UUID) ? " AND " + Grid::workspace_security_clause("rows") : "") +
          ((self.uuid == GRID_UUID) ? " AND " + Grid::grid_security_clause("rows") : "")
    Row.find_by_sql([sql])[0]
    rescue ActiveRecord::StatementInvalid => exception
      log_error "Grid#row_select_entity_by_uuid #{exception.to_s} #{sql} [#{to_s}]"
      nil
  end

  # Selects a row based on its uuid or uri.
  def row_select_entity_by_uri(uri, version=nil)
    row_select_entity_by_uuid(nil, uri, version)
  end

  # Selects a row based on its uuid and version number.
  def row_select_entity_by_uuid_version(uuid, version=nil)
    row_select_entity_by_uuid(uuid, nil, version)
  end

  def row_all_versions(uuid)
    log_debug "Grid#row_all_versions(#{uuid}) [#{to_s}]"
    if not @can_select_data
      log_security_warning "Grid#row_all_versions Can't select data"
      return []
    end
    Row.find_by_sql(["SELECT #{row_all_select_columns}" + 
                     " FROM grids grids, #{@db_table} rows" + 
                     (@has_translation ? ", #{@db_loc_table} row_locs" : "") +
                     " WHERE grids.uuid = #{quote(self.uuid)}" +
                     " AND " + as_of_date_clause("grids") +
                     " AND " + Grid::grid_security_clause("grids") + 
                     " AND rows.uuid = :uuid" +
                     (@has_mapping ? "" : " AND rows.grid_uuid = :grid_uuid") +
                     (@has_translation ? 
                        " AND rows.uuid = row_locs.uuid" +
                        " AND rows.version = row_locs.version" +
                        " AND " + locale_clause("row_locs") : "") + 
                     ((self.uuid == WORKSPACE_UUID) ? 
                        " AND " + Grid::workspace_security_clause("rows") : "") + 
                     ((self.uuid == GRID_UUID) ? 
                         " AND " + Grid::grid_security_clause("rows") : "") + 
                     " ORDER BY rows.begin", 
                       {:uuid => uuid, :grid_uuid => self.uuid}])
  end

  def row_locales(uuid, version, all=true)
    log_debug "Grid#row_locales(#{uuid}, #{version.to_s}) [#{to_s}]"
    if not @can_select_data
      log_security_warning "Grid#row_locales Can't select data"
      return []
    end
    RowLoc.find_by_sql(["SELECT #{row_loc_select_columns}" + 
                        " FROM grids grids, #{@db_loc_table} row_locs" +
                        " WHERE grids.uuid = #{quote(self.uuid)}" +
                        " AND " + as_of_date_clause("grids") +
                        " AND " + Grid::grid_security_clause("grids") + 
                        " AND row_locs.uuid = :uuid" +
                        " AND row_locs.version = :version" +
                        (all ? "" : " AND row_locs.base_locale = row_locs.locale") +
                        " ORDER BY row_locs.locale", 
                         {:uuid => uuid, :version => version}])
  end

  def row_begin_duplicate_exists?(row, begin_date)
    log_debug "Grid#row_begin_duplicate_exists?(begin_date=#{begin_date})"
    sql = "SELECT id" + 
          " FROM #{@db_table}" + 
          " WHERE uuid = #{quote(row.uuid)}" +
          (@has_mapping ? "" : " AND grid_uuid = #{quote(self.uuid)}") +
          " AND begin = #{quote(begin_date)}" +
          (row.id.present? ? " AND id != #{quote(row.id)}" : "")
    connection.select_value(sql).present?
  end

  def row_enabled_version_exists?(row, new_version)
    log_debug "Grid#row_enabled_version_exists?(" + 
              "new_version=#{new_version.to_s})"
    sql = "SELECT id" + 
          " FROM #{@db_table}" + 
          " WHERE uuid = #{quote(row.uuid)}" +
          (@has_mapping ? "" : " AND grid_uuid = #{quote(self.uuid)}") +
          " AND enabled = #{quote(true)}" +
          (row.id.present? ? " AND id != #{quote(row.id)}" : "")
    connection.select_value(sql).present?
  end

  # Selects greatest version number of data based on uuid
  def row_max_version(uuid)
    log_debug "Grid#row_max_version(uuid=#{uuid}) [#{to_s}]"
    sql = "SELECT max(version)" + 
          " FROM #{@db_table}" + 
          " WHERE uuid = #{quote(uuid)}" +
          (@has_mapping ? "" : " AND grid_uuid = #{quote(self.uuid)}")
    connection.select_value(sql).to_i
  end

  def row_title(row, full=false)
    summary = (row.present? and self.has_name) ? row.name : ""
    if not full and row.present?
      return summary if summary.length > 0
      count = 0
      load if not loaded
      @filtered_columns.each do |column|
        if [COLUMN_TYPE_REFERENCE, 
            COLUMN_TYPE_STRING, 
            COLUMN_TYPE_TEXT].include?(column.kind)
          value = row.read_referenced_name(column)
          if value.length > 0
            summary = summary + " | " if count > 0
            summary = summary + value
            count += 1
          end
        end
      end
      return summary
    end
    return summary.length > 0 ? summary : "[-]" if full 
    I18n.t('general.undefined')
  end

  def row_loc_select_entity_by_uuid(uuid, version=0)
    log_debug "Grid#row_loc_select_entity_by_uuid(uuid=#{uuid}, " +
              "version=#{version}) [#{to_s}]"
    if not @can_select_data
      log_security_warning "Grid#row_loc_select_entity_by_uuid Can't select data"
      return []
    end
    RowLoc.find_by_sql(["SELECT #{row_loc_select_columns}" + 
                        " FROM grids grids, #{@db_loc_table} row_locs" + 
                        " WHERE grids.uuid = #{quote(self.uuid)}" +
                        " AND " + as_of_date_clause("grids") +
                        " AND " + Grid::grid_security_clause("grids") + 
                        " AND row_locs.uuid = :uuid" +
                        (version > 0 ? 
                          " AND row_locs.version = :version" : 
                          ""), 
                          {:version => version, :uuid => uuid}])
  end

  def row_initialization(row, filters)
    log_debug "Grid#row_initialization(filters=#{filters.inspect}) [#{to_s}]"
    row.initialization
    column_all.each do |column|
      log_debug "Grid#default_row_value column=#{column.to_s}"
      initialized = false
      if filters.present?
        filters.each do |filter|
          if column.uuid == filter[:column_uuid]
            row.write_value(column, filter[:row_uuid])
            initialized = true
          end
        end
      end
      if not initialized
        row.write_value(column, nil)
        if self.uuid == COLUMN_UUID and column.uuid == COLUMN_KIND_UUID
          row.write_value(column, COLUMN_TYPE_STRING)
        elsif self.uuid == COLUMN_UUID and column.uuid == COLUMN_DISPLAY_UUID
          row.write_value(column, 1)
        elsif self.uuid == GRID_UUID and column.uuid == GRID_HAS_NAME_UUID
          row.write_value(column, true)
        elsif self.uuid == GRID_UUID and column.uuid == GRID_HAS_DESCRIPTION_UUID
          row.write_value(column, true)
        end
      end
    end
  end

  # Creates a new associated locale row.
  def new_loc
    loc = grid_locs.new
    loc.uuid = self.uuid
    loc.version = self.version
    loc
  end

  # Exports the entity into an .xml output.
  def export(xml)
    log_debug "Grid#export [#{to_s}]"
    xml.grid(:title => self.name) do
      super(xml)
      xml.workspace_uuid(self.workspace_uuid, :title => self.workspace_name)
      xml.has_name(self.has_name) if self.has_name
      xml.has_description(self.has_description) if self.has_description
      Grid.locales(grid_locs, self.uuid, self.version).each do |loc|
        if loc.base_locale == loc.locale
          log_debug "Grid#export locale #{loc.base_locale}"
          xml.locale do
            xml.base_locale(loc.base_locale)
            xml.locale(loc.locale)
            xml.name(loc.name)
            xml.description(loc.description) if loc.description.present?
          end
        end
      end
    end
  end

  # Copies attributes from the object to the target entity.
  def copy_attributes(entity)
    log_debug "Grid#copy_attributes"
    super
    entity.workspace_uuid = self.workspace_uuid
    entity.has_name = self.has_name
    entity.has_description = self.has_description
  end

  # Imports attribute value from the xml flow into the object.
  def import_attribute(xml_attribute, xml_value)
    log_debug "Grid#import_attribute(#{xml_attribute}, #{xml_value})"
    case xml_attribute
      when GRID_WORKSPACE_UUID then self.workspace_uuid = xml_value
      when GRID_HAS_NAME_UUID then self.has_name = ['true','t','1'].include?(xml_value)
      when GRID_HAS_DESCRIPTION_UUID then self.has_description = ['true','t','1'].include?(xml_value)
      when GRID_URI_UUID then self.uri = xml_value
    end
  end

  # Imports the instance of the object in the database,
  # as a new instance or as an update of an existing instance.
  def import!
    log_debug "Grid#import!"
    grid = Grid.select_entity_by_uuid_version(Grid, self.uuid, self.version)
    if grid.present?
      if self.revision > grid.revision 
        log_debug "Grid#import! update"
        copy_attributes(grid)
        grid.update_user_uuid = Entity.session_user_uuid
        grid.updated_at = Time.now
        make_audit(Audit::IMPORT)
        grid.save!
        grid.update_dates!(Grid)
        return "updated"
      else
        log_debug "Grid#import! skip update"
        return "skipped"
      end
    else
      log_debug "Grid#import! new"
      self.create_user_uuid = self.update_user_uuid = Entity.session_user_uuid
      self.created_at = self.updated_at = Time.now
      make_audit(Audit::IMPORT)
      save!
      update_dates!(Grid)
      return "inserted"
    end
    ""
  end

  # Imports the given loc data into the appropriate locale row.
  # Fetches on the collection of local row and copies the attributes
  # of the provided loc data into the row that matches the language.
  def import_loc!(loc)
    log_debug "Grid#import_loc!"
    import_loc_base!(Grid.locales(grid_locs, self.uuid, self.version), loc)
  end

  # Creates local row for all the installed languages
  # that is not created yet for the given collection.
  # This insures on row exists for any installed language.
  def create_missing_loc!
    log_debug "Grid#create_missing_loc!"
    create_missing_loc_base!(Grid.locales(grid_locs, self.uuid, self.version))
  end

  # Exports row in .xml format.
  # For rows corresponding to workspaces, grids or columns, associated rows
  # are also exported, so a grid is exported with columns, a workspace is associated
  # with grids and so on. 
  def row_export(xml, row)
    log_debug "Grid#row_export(row=#{row}) [#{to_s}]"
    if row.present?
      xml.row(:title => row_title(row), :grid_uuid => self.uuid, :grid => to_s) do
        row.export(xml)
        xml.grid_uuid(self.uuid, :title => name)
        column_all.each do |column|
          xml.data(row.read_value(column), :uuid => column.uuid, :name => column.name)
        end
        row_locales(row.uuid, row.version).each do |loc|
          if loc.base_locale == loc.locale
            xml.locale do
              xml.base_locale(loc.base_locale)
              xml.locale(loc.locale)
              xml.name(loc.name) if loc.name.present?
              xml.description(loc.description) if loc.description.present?
            end
          end
        end
      end
      if self.uuid == WORKSPACE_UUID
        grid_def = Grid::select_entity_by_uuid(Grid, GRID_UUID)
        if grid_def.present?
          grid_def.load
          log_debug "Grid#row_export select grids in the workspace"
          rows = grid_def.row_all([{:column_uuid => GRID_WORKSPACE_UUID,
                                    :row_uuid => row.uuid}], nil, -1, true)
          for child in rows
            grid_def.row_export(xml, child)
          end
        end
      end
      if self.uuid == GRID_UUID
        grid_def = Grid::select_entity_by_uuid(Grid, GRID_MAPPING_UUID)
        if grid_def.present?
          grid_def.load
          log_debug "Grid#row_export select mapping in the data grid"
          rows = grid_def.row_all([{:column_uuid => GRID_MAPPING_GRID_UUID,
                                    :row_uuid => row.uuid}], nil, -1, true)
          for child in rows
            grid_def.row_export(xml, child)
          end
        end
        grid_def = Grid::select_entity_by_uuid(Grid, COLUMN_UUID)
        if grid_def.present?
          grid_def.load
          log_debug "Grid#row_export select columns in the data grid"
          rows = grid_def.row_all([{:column_uuid => COLUMN_GRID_UUID,
                                    :row_uuid => row.uuid}], nil, -1, true)
          for child in rows
            grid_def.row_export(xml, child)
          end
        end
      end
      if self.uuid == COLUMN_UUID
        grid_def = Grid::select_entity_by_uuid(Grid, COLUMN_MAPPING_UUID)
        if grid_def.present?
          grid_def.load
          log_debug "Grid#row_export select mapping in the column"
          rows = grid_def.row_all([{:column_uuid => COLUMN_MAPPING_COLUMN_UUID,
                                    :row_uuid => row.uuid}], nil, -1, true)
          for child in rows
            grid_def.row_export(xml, child)
          end
        end
      end
    end
  end

  # Indicates if the user can create a row into the grid.
  # Authorization depends on the workspace security settings
  # the grid is attached to, but also to the workspace security settings
  # the row will be attached to, based on the filters.
  # Thus, if the filter contains a reference to a workspace, 
  # the security applies based on the workspace provided by the filter.
  def can_create_row?(filters=nil)
    return false if not(@can_create_data or [WORKSPACE_SHARING_UUID,
                                             GRID_UUID,
                                             COLUMN_UUID].include?(self.uuid))
    if filters.present?
      filters.each do |filter|
        column_uuid = filter[:column_uuid]
        row_uuid = filter[:row_uuid]
        column_all.each do |column|
          if column.uuid == column_uuid and 
             column.kind == COLUMN_TYPE_REFERENCE and
             column.grid_reference_uuid == WORKSPACE_UUID
            security = get_security_workspace(row_uuid)
            return [ROLE_READ_WRITE_ALL_UUID,
                    ROLE_TOTAL_CONTROL_UUID].include?(security) if not security.nil?
          elsif column.uuid == column_uuid and 
                column.kind == COLUMN_TYPE_REFERENCE and
                column.grid_reference_uuid == GRID_UUID
            grid = Grid::select_entity_by_uuid(Grid, row_uuid)
            if grid.present?
              security = get_security_workspace(grid.workspace_uuid)
              return [ROLE_READ_WRITE_ALL_UUID,
                      ROLE_TOTAL_CONTROL_UUID].include?(security) if not security.nil?
            end
          end
        end
      end
    end
    return @can_create_data
  end

  # Indicates if the user can update a specific row.
  # Authorization depends on the workspace security settings
  # the row is attached to. The way the row is attached to a workspace depends
  # on the grid. The row could be a grid itself, a workspace itself or a column.
  def can_update_row?(row)
    return false if not(@can_update_data or [WORKSPACE_UUID,
                                             WORKSPACE_SHARING_UUID,
                                             GRID_UUID,
                                             COLUMN_UUID].include?(self.uuid))
    security = nil
    if self.uuid == WORKSPACE_UUID
      security = get_security_workspace(row.uuid)
    elsif self.uuid == GRID_UUID or self.uuid == WORKSPACE_SHARING_UUID
      security = get_security_workspace(row.workspace_uuid)
    elsif self.uuid == COLUMN_UUID and row.grid.present?
      security = get_security_workspace(row.grid.workspace_uuid)
    end
    return [ROLE_READ_WRITE_ALL_UUID,
            ROLE_TOTAL_CONTROL_UUID].include?(security) if security.present?
    return (@can_update_data or (row.create_user_uuid == Entity.session_user_uuid))
  end

  def create_row!(row, filters=nil)
    log_debug "Grid#create_row!(row=#{row.inspect})"
    if not can_create_row?(filters)
      log_security_warning "Grid#create_row! Can't create data"
      row.errors.add(:uuid, I18n.t('error.cant_create'))
      raise "Grid#create_row! Can't create data"
    end
    row.created_at = row.updated_at = Time.now
    row.create_user_uuid = row.update_user_uuid = Entity.session_user_uuid
    sql = "INSERT INTO #{@db_table}" +
          "(#{row_all_insert_columns})" +
          " VALUES(#{row_all_insert_values(row)})"
    self.id = connection.insert(sql, 
                                "#{self.class.name} Create",
                                self.class.primary_key, 
                                self.id, 
                                self.class.sequence_name)
    row.make_audit(Audit::CREATE)
    true
  end

  def create_row_loc!(row, row_loc, filters=nil)
    log_debug "Grid#create_row_loc!(row_loc=#{row_loc.inspect})"
    if not can_create_row?(filters)
      log_security_warning "Grid#create_row_loc! Can't create data"
      row.errors.add(:uuid, I18n.t('error.cant_create'))
      raise "Grid#create_row_loc! Can't create data"
    end
    sql = "INSERT INTO #{@db_loc_table}" +
          "(#{row_all_insert_loc_columns})" +
          " VALUES(#{row_all_insert_loc_values(row_loc)})"
    self.id = connection.insert(sql, 
                                "#{self.class.name} Create",
                                self.class.primary_key, 
                                self.id, 
                                self.class.sequence_name)
    true
  end

  def update_row!(row, audit=Audit::UPDATE)
    log_debug "Grid#update_row!(row=#{row.inspect})"
    if not can_update_row?(row)
      log_security_warning "Grid#update_row! Can't update data"
      row.errors.add(:uuid, I18n.t('error.cant_update'))
      raise "Grid#update_row! Can't update data"
    end
    row.updated_at = Time.now
    row.update_user_uuid = Entity.session_user_uuid
    sql = "UPDATE #{@db_table}" +
          " SET #{row_all_update_values(row)}" +
          " WHERE id = #{quote(row.id)}" +
          " AND lock_version = #{quote(row.lock_version)}"
    connection.update(sql, "#{self.class.name} Update")
    row.lock_version += 1
    row.make_audit(audit)
    true
  end

  def update_row_loc!(row, row_loc)
    log_debug "Grid#update_row_loc!(row_loc=#{row_loc.inspect})"
    if not can_update_row?(row)
      log_security_warning "Grid#update_row_loc! Can't update data"
      row.errors.add(:uuid, I18n.t('error.cant_update'))
      raise "Grid#update_row_loc! Can't update data"
    end
    sql = "UPDATE #{@db_loc_table}" +
          " SET #{row_loc_update_values(row_loc)}" +
          " WHERE id = #{quote(row_loc.id)}" +
          " AND lock_version = #{quote(row_loc.lock_version)}"
    connection.update(sql, "#{self.class.name} Update")
    true
  end

  def row_update_dates!(uuid)
    log_debug "Grid#row_update_dates!(uuid=#{uuid}"
    previous_item = nil
    last_item = nil
    row_all_versions(uuid).each do |item|
      if item.enabled
        if previous_item.present? and previous_item.end != item.begin-1
          log_debug "Grid#row_update_dates! previous_item set end date"
          previous_item.end = item.begin-1 if item.begin > Entity.begin_of_time
          previous_item.update_user_uuid = Entity.session_user_uuid
          update_row!(previous_item)
        end
        previous_item = item
      end
      last_item = item
    end
    log_debug "Grid#row_update_dates! last_item=#{last_item.inspect}"
    if last_item.present? and last_item.end != @@end_of_time
      log_debug "Entity#update_dates last_item set end date"
      last_item.end = @@end_of_time
      last_item.update_user_uuid = Entity.session_user_uuid
      update_row!(last_item)
    end
  end

  def mapping_select_entity_by_uuid_version(uuid, version)
    log_debug "Grid#mapping_select_entity_by_uuid_version(uuid=#{uuid}, version=#{version})"
    grid_mappings.find(:all, 
                       :conditions => 
                          ["uuid = :uuid and version = :version", 
                          {:uuid => uuid, :version => version}])[0]
  end

  # Validates row is correct according to grid columns definition.
  # Row is validated against filters also in order to control data
  # match with the scope of selected data.
  def row_validate(row, phase, filters)
    log_debug "Grid#row_validate(phase=#{phase}) [#{to_s}]"
    validated = true
    for column in column_all
      attribute = (phase == Grid::PHASE_CREATE) ? column.default_physical_column : column.physical_column
      value = row.read_value(column)
      log_debug "Grid#row_validate(phase=#{phase}) control #{column.name}: #{value}"
      if filters.present?
        filters.each do |filter|
          if column.uuid == filter[:column_uuid] and value != filter[:row_uuid]
            validated = false
            row.errors.add(attribute, I18n.t('error.cant_select', :column => column))
          end
        end
      end
      if column.required
        if value.blank?
          log_debug "Grid#row_validate(phase=#{phase}) required"
          validated = false
          row.errors.add(attribute, I18n.t('error.required', :column => column))
        end
      end
      if column.regex.present?
        if not Regexp.new(column.regex).match(value)
          log_debug "Grid#row_validate(phase=#{phase}) regex"
          validated = false
          row.errors.add(attribute, I18n.t('error.badformat', :column => column))
        end
      end
    end
    validated
  end

  # Validates locale row is correct.
  def row_loc_validate(row, row_loc, phase)
    log_debug "Grid#row_loc_validate(phase=#{phase}) [#{to_s}]"
    validated = true
    if self.has_name
      if row_loc.name.blank?
          validated = false
          row.errors.add(:name, I18n.t('error.mis_name'))
      end
    end
    validated
  end

private

  def quote(sql) ; connection.quote(sql) ; end

  def self.all_select_columns
    "grids.id, grids.uuid, grids.version, grids.lock_version, " + 
    "grids.begin, grids.end, grids.enabled, grids.workspace_uuid, " + 
    "grids.has_name, grids.has_description, grids.uri, " + 
    "grids.created_at, grids.updated_at, " +
    "grids.create_user_uuid, grids.update_user_uuid, " +
    "grid_locs.base_locale, grid_locs.locale, " +
    "grid_locs.name, grid_locs.description"
  end

  def column_all_select_columns
    "columns.id, columns.uuid, columns.uri, columns.version, columns.lock_version, " +
    "columns.begin, columns.end, columns.enabled, " +
    "columns.display, columns.kind, " + 
    "columns.grid_reference_uuid, " + 
    "columns.required, columns.regex, " + 
    "columns.created_at, columns.updated_at, " + 
    "columns.create_user_uuid, columns.update_user_uuid, " + 
    "column_locs.base_locale, column_locs.locale, " +
    "column_locs.name, column_locs.description"
  end

  def column_all_select_columns_mapping
    column_all_select_columns + ", column_mappings.db_column"
  end

  def row_all_select_columns
    columns = ""
    column_all.each{|column| columns << ", rows." + column.physical_column}
    "#{quote(self.uuid)} grid_uuid" +
    ", rows.id, rows.uuid, rows.uri, rows.version, rows.lock_version" + 
    ", rows.begin, rows.end, rows.enabled" + 
    ", rows.created_at, rows.updated_at" +
    ", rows.create_user_uuid, rows.update_user_uuid" +
    (@has_translation ? 
      ", row_locs.base_locale, row_locs.locale" + 
      ", row_locs.name, row_locs.description" : 
      "") + 
    columns
  end

  def row_loc_select_columns
    "row_locs.id, row_locs.uuid, " +
    "row_locs.version, row_locs.lock_version, " +
    "row_locs.base_locale, row_locs.locale, " + 
    "row_locs.name, row_locs.description" 
  end

  def row_all_insert_columns
    columns = ""
    column_all.each{|column| columns << ", " + column.physical_column}
    "uuid, uri, version, lock_version, begin, end, enabled" + 
    (@has_mapping ? "" : ", grid_uuid") + 
    ", created_at, updated_at, create_user_uuid, update_user_uuid" + 
    columns
  end

  def row_all_insert_loc_columns
    "uuid, version, lock_version, locale, base_locale, name, description" 
  end

  def row_all_insert_values(row)
    columns = ""
    column_all.each{|column| columns << ", #{quote(row.read_value(column))}"}
    "#{quote(row.uuid)}, #{quote(row.uri)}, #{quote(row.version)}, #{quote(row.lock_version)}" +
    ", #{quote(row.begin)}, #{quote(row.end)}, #{quote(row.enabled)}" +
    (@has_mapping ? "" : ", #{quote(uuid)}") + 
    ", #{quote(row.created_at)}, #{quote(row.updated_at)}" +
    ", #{quote(row.create_user_uuid)}, #{quote(row.update_user_uuid)}" +
    columns
  end

  def row_all_insert_loc_values(row)
    "#{quote(row.uuid)}" +
    ", #{quote(row.version)}, #{quote(row.lock_version)}" +
    ", #{quote(row.locale)}, #{quote(row.base_locale)}" +
    ", #{quote(row.name)}, #{quote(row.description)}"
  end

  def row_all_update_values(row)
    columns = ""
    column_all.each{|column| columns << ", #{column.physical_column}=#{quote(row.read_value(column))}"}
    "lock_version=#{quote(row.lock_version+1)}" +
    ", uri=#{quote(row.uri)}" +
    ", begin=#{quote(row.begin)}, end=#{quote(row.end)}, enabled=#{quote(row.enabled)}" +
    ", updated_at=#{quote(row.updated_at)}, update_user_uuid=#{quote(row.update_user_uuid)}" +
    columns
  end

  def row_loc_update_values(row)
    "lock_version=#{quote(row.lock_version+1)}, base_locale=#{quote(row.base_locale)}" +
    ", name=#{quote(row.name)}, description=#{quote(row.description)}"
  end

  def load_workspace
    log_debug "Grid#load_workspace [#{to_s}]"
    self.workspace = Workspace.select_entity_by_uuid(Workspace, self.workspace_uuid)
    log_debug "Grid#load_workspace workspace=#{self.workspace.to_s}"
    if self.workspace.present? and self.workspace.default_role_uuid.present?
      @can_select_data = true
      @can_create_data = [ROLE_READ_WRITE_UUID,
                          ROLE_READ_WRITE_ALL_UUID,
                          ROLE_TOTAL_CONTROL_UUID].include?(self.workspace.default_role_uuid)
      @can_update_data = [ROLE_READ_WRITE_UUID,
                          ROLE_READ_WRITE_ALL_UUID,
                          ROLE_TOTAL_CONTROL_UUID].include?(self.workspace.default_role_uuid)
    end
    sql = "SELECT workspace_sharings.role_uuid" + 
          " FROM workspace_sharings" + 
          " WHERE workspace_sharings.workspace_uuid = '#{self.workspace_uuid}'" + 
          " AND workspace_sharings.user_uuid = '#{Entity.session_user_uuid}'" +
          " AND " + as_of_date_clause("workspace_sharings") +
          " LIMIT 1"
    security = Grid.find_by_sql([sql])[0]
    if security.present?
      @can_select_data = true
      @can_create_data = [ROLE_READ_WRITE_UUID,
                          ROLE_READ_WRITE_ALL_UUID,
                          ROLE_TOTAL_CONTROL_UUID].include?(security.role_uuid)
      @can_update_data = [ROLE_READ_WRITE_UUID,
                          ROLE_READ_WRITE_ALL_UUID,
                          ROLE_TOTAL_CONTROL_UUID].include?(security.role_uuid)
    end
    log_debug "Grid#load_workspace " +
              "@can_select_data=#{@can_select_data}," +
              "@can_create_data=#{@can_create_data}," +
              "@can_update_data=#{@can_update_data} [#{to_s}]"
  end

  # Loads in memory information about database mapping.
  def load_mapping
    log_debug "Grid#load_mapping [#{to_s}]"
    grid_mapping = grid_mappings.find(:first, 
                                      :select => "grid_mappings.db_table, grid_mappings.db_loc_table",
                                      :conditions => [as_of_date_clause("grid_mappings")])
    if grid_mapping.present?
      @db_table = grid_mapping.db_table if grid_mapping.db_table.present?
      @db_loc_table = grid_mapping.db_loc_table if grid_mapping.db_loc_table.present?
    else
      @db_table = "rows"
      @db_loc_table = "row_locs"
    end
    @has_mapping = grid_mapping.present?
    @has_translation = (@db_loc_table.present? and (self.has_name or self.has_description))
    log_debug "Grid#load_mapping self.has_name=#{self.has_name}, self.has_description=#{self.has_description}, @has_mapping=#{@has_mapping}, @has_translation=#{@has_translation} [#{to_s}]"
  end

  # Loads in memory information about columns.
  def load_columns(filters, skip_reference=false, skip_mapping=false)
    log_debug "Grid#load_columns " +
              "filters=#{filters},"+
              "skip_reference=#{skip_reference}," +
              "skip_mapping=#{skip_mapping}) [#{to_s}]"
    if not(skip_mapping) and @has_mapping
      @column_all = columns.find(:all,
                                 :joins => [:column_locs, :column_mappings],
                                 :select => column_all_select_columns_mapping,
                                 :conditions =>
                                     [as_of_date_clause("columns") +
                                      " AND column_locs.version = columns.version" +
                                      " AND " + locale_clause("column_locs")],
                                 :order => "columns.id")
    else
      @column_all = columns.find(:all,
                                 :joins => :column_locs,
                                 :select => column_all_select_columns,
                                 :conditions =>
                                     [as_of_date_clause("columns") +
                                      " AND column_locs.version = columns.version" +
                                      " AND " + locale_clause("column_locs")],
                                 :order => "columns.id")
    end
    @filtered_columns = Array.new(@column_all)
    index_reference = index_date = index_integer = index_decimal = index_string = 0
    column_all.each do |column|
      unless column.loaded
        case column.kind
          when COLUMN_TYPE_REFERENCE then
            index_reference += 1 
            number = index_reference
          when COLUMN_TYPE_DATE then 
            index_date += 1 
            number = index_date
          when COLUMN_TYPE_INTEGER then 
            index_integer += 1 
            number = index_integer
          when COLUMN_TYPE_DECIMAL then 
            index_decimal += 1 
            number = index_decimal
          else 
            index_string += 1 
            number = index_string
        end
        if column.display.blank? or column.display == 0
          @filtered_columns.delete(column)
        else
          column.load_cached_information(self,
                                         number,
                                         skip_reference,
                                         (not(@has_mapping) or skip_mapping),
                                         (not(skip_mapping) and @has_mapping) ? column.db_column : nil)
        end
        filters.each{|filter| @filtered_columns.delete(column) if filter[:column_uuid] == column.uuid} unless filters.nil?
      end
    end
    @filtered_columns.sort!{|a, b| a.display <=> b.display}
  end

  def get_security_workspace(uuid)
    log_debug "Grid#get_security_workspace [#{to_s}]"
    sql = "SELECT workspace_sharings.role_uuid" + 
          " FROM workspace_sharings" + 
          " WHERE workspace_sharings.workspace_uuid = '#{uuid}'" + 
          " AND workspace_sharings.user_uuid = '#{Entity.session_user_uuid}'" +
          " AND " + as_of_date_clause("workspace_sharings") +
          " LIMIT 1"
    security = Grid.find_by_sql([sql])[0]
    return security.role_uuid if not security.nil?
    sql = "SELECT default_role_uuid, create_user_uuid" + 
          " FROM workspaces workspace_security" + 
          " WHERE workspace_security.uuid = '#{uuid}'" + 
          " AND workspace_security.default_role_uuid is not null" +
          " AND " + as_of_date_clause("workspace_security") +
          " LIMIT 1"
    security = Grid.find_by_sql([sql])[0]
    security.nil? ? nil : security.default_role_uuid
  end
end