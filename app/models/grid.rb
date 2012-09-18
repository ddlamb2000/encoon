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
  ROOT_UUID = 'fe92c430-ea30-012c-1057-00166f92f624'
  ROOT_WORKSPACE_UUID = '1fe00120-ea31-012c-1065-00166f92f624'
  ROOT_HAS_NAME_UUID = '78a7e2d1-293a-012d-2869-4417fe7fde95'
  ROOT_HAS_DESCRIPTION_UUID = '8afbf6b1-293a-012d-1701-4417fe7fde95'
  
  HOME_GRID_UUID = 'cf01e9a0-d59d-012f-590e-4417fe7fde95'
  HOME_ROW_UUID = 'f38b12c0-d59d-012f-590e-4417fe7fde95'

  PHASE_CREATE = 'create'
  PHASE_NEW_VERSION = 'new_version'
  PHASE_UPDATE = 'update'
  
  DISPLAY_ROWS_LIMIT = 10
  DISPLAY_ROWS_LIMIT_FULL = 50
  
  belongs_to :workspace, :foreign_key => "workspace_uuid", :primary_key => "uuid"
  has_many :grid_locs, :foreign_key => "uuid", :primary_key => "uuid"
  has_many :columns,  :foreign_key => "grid_uuid", :primary_key => "uuid"
  has_many :column_locs,  :foreign_key => "uuid", :primary_key => "uuid", :through => :columns
  has_many :rows,  :foreign_key => "grid_uuid", :primary_key => "uuid"
  has_many :grid_mappings,  :foreign_key => "grid_uuid", :primary_key => "uuid"
  validates_presence_of :workspace_uuid
  validates_associated :workspace
  
  @can_select_data = false
  @can_create_data = false
  @can_update_data = false

  # Loads in memory the structure of the data grid
  # and all the information to render the grid (columns
  # definition, table mapping).
  def load_cached_grid_structure(filters=nil, skip_mapping=false)
    log_debug "Grid#load_grid_structure(" +
              "filters=#{filters.inspect},"+
              "skip_mapping=#{skip_mapping}) [#{to_s}]"
    load_workspace
    load_cached_mapping
    load_columns(filters, false, skip_mapping)
    load_security
  end
  
  def load_cached_grid_structure_reference
    log_debug "Grid#load_cached_grid_structure_reference [#{to_s}]"
    load_workspace
    load_cached_mapping
    load_columns(nil, true)
    load_security
  end
  
  def is_preloaded? ; @columns.present? ; end
  def has_name? ; self.has_name ; end
  def has_description? ; self.has_description ; end
  def has_translation? ; @db_loc_table.present? and (has_name? or has_description?) ; end
  def has_mapping? ; @has_mapping ; end
  def can_select_data? ; @can_select_data ; end
  def can_update_data? ; @can_update_data ; end

  def can_create_data?(filters=nil)
    if @can_create_data
      if filters.present?
        filters.each do |filter|
          column_uuid = filter[:column_uuid]
          row_uuid = filter[:row_uuid]
          column_all.each do |column|
            if column.uuid == column_uuid and 
                column.kind == Column::REFERENCE and
                column.grid_reference_uuid == Workspace::ROOT_UUID
              security = load_security_workspace(row_uuid)
              return Role::ROLE_TOTAL_CONTROL_UUID == security if not security.nil?
            end
          end
        end
      end
      return true
    end
    false
  end

  def can_update?(row)
    if can_update_data?
      return true if row.create_user_uuid == Entity.session_user_uuid
      if self.uuid == Workspace::ROOT_UUID
        security = load_security_workspace(row.uuid)
        return Role::ROLE_TOTAL_CONTROL_UUID == security if not security.nil?
      elsif self.uuid == Grid::ROOT_UUID or self.uuid == WorkspaceSharing::ROOT_UUID
        security = load_security_workspace(row.workspace_uuid)
        return Role::ROLE_TOTAL_CONTROL_UUID == security if not security.nil?
      elsif self.uuid == Column::ROOT_UUID and row.grid.present?
        security = load_security_workspace(row.grid.workspace_uuid)
        return Role::ROLE_TOTAL_CONTROL_UUID == security if not security.nil?
      end
    end
    false    
  end

  # Selects data based on uuid.
  def self.select_entity_by_uuid(collection, uuid)
    log_debug "Grid#select_entity_by_uuid(" +
              "collection=#{collection}, " +
              "uuid=#{uuid})"
    collection.find(:first, 
                    :joins => :grid_locs,
                    :select => self.all_select_columns,
                    :conditions => 
                        ["grids.uuid = :uuid" +
                         " AND grid_locs.version = grids.version" +
                         " AND " + as_of_date_clause("grids") +
                         " AND " + locale_clause("grid_locs") +
                         " AND " + grid_security_clause("grids"), 
                        {:uuid => uuid}])
  end
  
  # Selects data based on workspace and uri.
  def self.select_entity_by_workspace_and_uri(collection, workspace_uuid, uri)
    log_debug "Grid#select_entity_by_workspace_and_uri(" +
              "collection=#{collection}, " +
              "uri=#{uri})"
    collection.find(:first, 
                    :joins => :grid_locs,
                    :select => self.all_select_columns,
                    :conditions => 
                        ["grid.workspace_uuid = :workspace_uuid" +
                         " AND grids.uri = :uri" +
                         " AND grid_locs.version = grids.version" +
                         " AND " + as_of_date_clause("grids") +
                         " AND " + locale_clause("grid_locs") +
                         " AND " + grid_security_clause("grids"), 
                        {:workspace_uuid => workspace_uuid,
                         :uri => uri}])
  end
  
  # Selects data based on id
  def self.select_entity_by_id(collection, id)
    log_debug "Grid#select_entity_by_id(" +
              "collection=#{collection}, " +
              "id=#{id})"
    collection.find(:first, 
                    :joins => :grid_locs,
                    :select => self.all_select_columns,
                    :conditions => 
                        ["grids.id = :id " +
                         " AND grid_locs.version = grids.version" +
                         " AND " + locale_clause("grid_locs") +
                         " AND " + grid_security_clause("grids"), 
                        {:id => id}]) 
  end
  
  def self.select_entity_by_uuid_version(collection, uuid, version)
    log_debug "Grid#select_entity_by_uuid_version(" +
              "collection=#{collection}, " +
              "uuid=#{uuid}, version=#{version})"
    collection.find(:first, 
                    :joins => :grid_locs,
                    :select => self.all_select_columns,
                    :conditions => 
                        ["grids.uuid = :uuid " + 
                         " AND grids.version = :version " + 
                         " AND grid_locs.version = grids.version " +
                         " AND " + locale_clause("grid_locs") + 
                         " AND " + grid_security_clause("grids"), 
                        {:uuid => uuid, 
                         :version => version}]) 
  end
  
  # Returns the name of the grid used as a reference
  def reference_name
    name + (workspace.present? ? " [" + workspace_name + "]" : "")
  end

  # Selects the reference rows attached to one grid
  # This is used to display a drop-down list
  def self.select_reference_rows(grid_uuid)
    log_debug "Grid#select_reference_rows(grid_uuid=#{grid_uuid})"
    grid = Grid.select_entity_by_uuid(Grid, grid_uuid)
    unless grid.nil?
      grid.load_cached_grid_structure
      grid.row_all(nil, nil, -1)
    end
  end

  # Selects the name of one reference row attached to one grid
  def select_reference_row_name(row_uuid)
    log_debug "Grid#select_reference_row_name(row_uuid=#{row_uuid})"
    row = row_select_entity_by_uuid(row_uuid)
    row.present? ? row_title(row) : ""
  end
  
  # Selects the description of one reference row attached to one grid
  def select_reference_row_description(row_uuid)
    log_debug "Grid#select_reference_row_description(row_uuid=#{row_uuid})"
    row = row_select_entity_by_uuid(row_uuid)
    row.present? ? row.description : ""
  end
  
  def workspace_name
    workspace = Workspace.select_entity_by_uuid(Workspace, self.workspace_uuid)
    workspace.present? ? workspace.name : ""
  end

  def select_grid_cast(row_uuid)
    log_debug "Grid#select_grid_cast(row_uuid=#{row_uuid}) [#{to_s}]"
    if self.uuid == ROOT_UUID
      grid = Grid::select_entity_by_uuid(Grid, row_uuid)
      grid.load_cached_grid_structure if grid.present?
      grid
    end
  end

  def filter_column_uuid
    attribute_present?(:filter_column_uuid) ? read_attribute(:filter_column_uuid) : nil
  end

  def filter_column_name
    attribute_present?(:filter_column_name) ? read_attribute(:filter_column_name) : nil
  end
  
  def get_filters_on_row(row)
    unless self.filter_column_uuid.nil?
      [{:column_uuid => self.filter_column_uuid,
        :row_uuid => row.uuid,
        :row_name => row.name}]
    else
      []
    end
  end
  
  # Selects the grids attached to one row via one or more columns
  def self.select_referenced_grids(uuid)
    log_debug "Grid#select_referenced_grids(uuid=#{uuid}) [#{to_s}]"
    Grid.find_by_sql(["SELECT grids.id," + 
                      " grids.uuid," + 
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
                       {:kind => Column::REFERENCE, 
                        :uuid => uuid}])
  end
  
  def column_all
    log_error "Grid#column_all Data grid structure " + 
              "#{self.uuid} isn't preloaded!!!" if @all_columns.nil?
    @all_columns
  end

  def filtered_columns
    log_error "Grid#filtered_columns Data grid structure " + 
              "#{self.uuid} isn't preloaded!!!" if @columns.nil?
    @columns
  end

  def column_select_by_uuid(uuid)
    for column in self.column_all
      return column if column.uuid == uuid
    end
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
    log_debug "Grid#column_select_entity_by_uuid_version(uuid=#{uuid}," + 
              "version=#{version}) [#{to_s}]"
    columns.find(:first, 
                 :joins => :column_locs,
                 :select => column_all_select_columns,
                 :conditions => 
                      ["columns.uuid = :uuid " + 
                       "AND columns.version = :version " + 
                       "AND column_locs.version = columns.version " +
                       "AND " + locale_clause("column_locs"), 
                      {:uuid => uuid, 
                       :version => version}]) 
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
          if column.uuid == column_uuid and column.kind == Column::REFERENCE
            log_debug "Grid#row_all found reference"
            conditions << " AND rows.#{column.physical_column} = #{quote(row_uuid)}"
          end
        end
      end
    end
    if search.present? and has_translation?
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
          (has_translation? ? ", #{@db_loc_table} row_locs" : "") +
          " WHERE grids.uuid = #{quote(self.uuid)}" +
          " AND " + as_of_date_clause("grids") +
          " AND " + Grid::grid_security_clause("grids") + 
          " AND " + as_of_date_clause("rows") +
          (has_mapping? ? "" : " AND rows.grid_uuid = #{quote(self.uuid)}") +
          (has_translation? ? 
              " AND rows.uuid = row_locs.uuid" +
              " AND rows.version = row_locs.version" +
              " AND " + locale_clause("row_locs") : "") +
          ((self.uuid == Workspace::ROOT_UUID) ? 
              " AND " + Grid::workspace_security_clause("rows") : "") + 
          ((self.uuid == Grid::ROOT_UUID) ? 
              " AND " + Grid::grid_security_clause("rows") : "") + 
          conditions +
          ((has_translation? and not count) ? " ORDER BY row_locs.name" : "") +
          (count ? "" : " LIMIT #{offset}, #{limit}") 
    count ? connection.select_value(sql).to_i : Row.find_by_sql([sql])
    rescue ActiveRecord::StatementInvalid => exception
      log_error "Grid#row_all #{exception.to_s}" +
                ", data grid='#{name}', sql=#{sql}"
      count ? 0 : []
  end

  # Selects data based on uuid
  def row_select_entity_by_uuid(uuid)
    log_debug "Grid#row_select_entity_by_uuid(uuid=#{uuid}) [#{to_s}]"
    if not can_select_data?
      log_security_warning "Grid#row_select_entity_by_uuid Can't select data"
      return nil
    end
    sql = "SELECT #{row_all_select_columns}" + 
          " FROM grids grids, #{@db_table} rows" + 
          (has_translation? ? ", #{@db_loc_table} row_locs" : "") +
          " WHERE grids.uuid = #{quote(self.uuid)}" +
          " AND " + as_of_date_clause("grids") +
          " AND " + Grid::grid_security_clause("grids") + 
          " AND rows.uuid = #{quote(uuid)}" +
          " AND " + as_of_date_clause("rows") +
          (has_mapping? ? "" : " AND rows.grid_uuid = #{quote(self.uuid)}") +
          (has_translation? ? 
              " AND rows.uuid = row_locs.uuid" +
              " AND rows.version = row_locs.version" +
              " AND " + locale_clause("row_locs") : "") +
          ((self.uuid == Workspace::ROOT_UUID) ? 
              " AND " + Grid::workspace_security_clause("rows") : "") + 
          ((self.uuid == Grid::ROOT_UUID) ? 
              " AND " + Grid::grid_security_clause("rows") : "") 
    Row.find_by_sql([sql])[0]
    rescue ActiveRecord::StatementInvalid => exception
      log_error "Grid#row_select_entity_by_uuid #{exception.to_s}" +
                ", data grid='#{name}', sql=#{sql}"
      nil
  end
  
  # Selects data based on uri.
  def row_select_entity_by_uri(uri)
    log_debug "Grid#row_select_entity_by_uri(uri=#{uri}) [#{to_s}]"
    if not can_select_data?
      log_security_warning "Grid#row_select_entity_by_uri Can't select data"
      return nil
    end
    sql = "SELECT #{row_all_select_columns}" + 
          " FROM grids grids, #{@db_table} rows" + 
          (has_translation? ? ", #{@db_loc_table} row_locs" : "") +
          " WHERE grids.uuid = #{quote(self.uuid)}" +
          " AND " + as_of_date_clause("grids") +
          " AND " + Grid::grid_security_clause("grids") + 
          " AND rows.uri = #{quote(uri)}" +
          " AND " + as_of_date_clause("rows") +
          (has_mapping? ? "" : " AND rows.grid_uuid = #{quote(self.uuid)}") +
          (has_translation? ? 
              " AND rows.uuid = row_locs.uuid" +
              " AND rows.version = row_locs.version" +
              " AND " + locale_clause("row_locs") : "") +
          ((self.uuid == Workspace::ROOT_UUID) ? 
              " AND " + Grid::workspace_security_clause("rows") : "") + 
          ((self.uuid == Grid::ROOT_UUID) ? 
              " AND " + Grid::grid_security_clause("rows") : "") 
    Row.find_by_sql([sql])[0]
    rescue ActiveRecord::StatementInvalid => exception
      log_error "Grid#row_select_entity_by_uri #{exception.to_s}" +
                ", data grid='#{name}', sql=#{sql}"
      nil
  end
  
  def row_select_entity_by_uuid_version(uuid, version)
    load_cached_grid_structure if not is_preloaded?
    log_debug "Grid#row_select_entity_by_uuid_version(uuid=#{uuid}, " + 
              "version=#{version}) [#{to_s}]"
    if not can_select_data?
      log_security_warning "Grid#row_select_entity_by_uuid_version Can't select data"
      return nil
    end
    Row.find_by_sql(["SELECT #{row_all_select_columns}" + 
                     " FROM grids grids, #{@db_table} rows" + 
                     (has_translation? ? ", #{@db_loc_table} row_locs" : "") +
                     " WHERE grids.uuid = #{quote(self.uuid)}" +
                     " AND " + as_of_date_clause("grids") +
                     " AND " + Grid::grid_security_clause("grids") + 
                     " AND rows.uuid = :uuid" +
                     " AND rows.version = :version" +
                     (has_mapping? ? "" : " AND rows.grid_uuid = :grid_uuid") +
                     (has_translation? ? 
                        " AND rows.uuid = row_locs.uuid" +
                        " AND rows.version = row_locs.version" +
                        " AND " + locale_clause("row_locs") : "") + 
                     ((self.uuid == Workspace::ROOT_UUID) ? 
                        " AND " + Grid::workspace_security_clause("rows") : "") + 
                     ((self.uuid == Grid::ROOT_UUID) ? 
                         " AND " + Grid::grid_security_clause("rows") : ""), 
                       {:grid_uuid => self.uuid,
                        :version => version,
                        :uuid => uuid}])[0]
    rescue ActiveRecord::StatementInvalid => exception
      log_error "Grid#row_select_entity_by_uuid_version #{exception.to_s}" +
                ", data grid='#{name}', sql=#{sql}"
      nil
  end
  
  # Selects data based on id
  def row_select_entity_by_id(id)
    log_debug "Grid#row_select_entity_by_id(id=#{id}) [#{to_s}]"
    if not can_select_data?
      log_security_warning "Grid#row_select_entity_by_id Can't select data"
      return nil
    end
    Row.find_by_sql(["SELECT #{row_all_select_columns}" + 
                     " FROM grids grids, #{@db_table} rows" + 
                     (has_translation? ? ", #{@db_loc_table} row_locs" : "") +
                     " WHERE grids.uuid = #{quote(self.uuid)}" +
                     " AND " + as_of_date_clause("grids") +
                     " AND " + Grid::grid_security_clause("grids") + 
                     " AND rows.id = :id" + 
                     (has_translation? ? 
                        " AND rows.uuid = row_locs.uuid" +
                        " AND rows.version = row_locs.version" +
                        " AND " + locale_clause("row_locs") : "") +
                     ((self.uuid == Workspace::ROOT_UUID) ? 
                        " AND " + Grid::workspace_security_clause("rows") : "") + 
                     ((self.uuid == Grid::ROOT_UUID) ? 
                         " AND " + Grid::grid_security_clause("rows") : ""), 
                       {:id => id}])[0]
    rescue ActiveRecord::StatementInvalid => exception
      log_error "Grid#row_select_entity_by_id #{exception.to_s}" +
                ", data grid='#{name}'"
      nil
  end
  
  def row_all_versions(uuid)
    log_debug "Grid#row_all_versions(uuid=#{uuid}) [#{to_s}]"
    if not can_select_data?
      log_security_warning "Grid#row_all_versions Can't select data"
      return []
    end
    Row.find_by_sql(["SELECT #{row_all_select_columns}" + 
                     " FROM grids grids, #{@db_table} rows" + 
                     (has_translation? ? ", #{@db_loc_table} row_locs" : "") +
                     " WHERE grids.uuid = #{quote(self.uuid)}" +
                     " AND " + as_of_date_clause("grids") +
                     " AND " + Grid::grid_security_clause("grids") + 
                     " AND rows.uuid = :uuid" +
                     (has_mapping? ? "" : " AND rows.grid_uuid = :grid_uuid") +
                     (has_translation? ? 
                        " AND rows.uuid = row_locs.uuid" +
                        " AND rows.version = row_locs.version" +
                        " AND " + locale_clause("row_locs") : "") + 
                     ((self.uuid == Workspace::ROOT_UUID) ? 
                        " AND " + Grid::workspace_security_clause("rows") : "") + 
                     ((self.uuid == Grid::ROOT_UUID) ? 
                         " AND " + Grid::grid_security_clause("rows") : "") + 
                     " ORDER BY rows.begin", 
                       {:uuid => uuid,
                        :grid_uuid => self.uuid}])
  end
  
  def row_all_locales(uuid, version)
    log_debug "Grid#row_all_locales(uuid=#{uuid}, " + 
              "version=#{version.to_s}) [#{to_s}]"
    if not can_select_data?
      log_security_warning "Grid#row_all_locales Can't select data"
      return []
    end
    RowLoc.find_by_sql(["SELECT #{row_loc_select_columns}" + 
                        " FROM grids grids, #{@db_loc_table} row_locs" +
                        " WHERE grids.uuid = #{quote(self.uuid)}" +
                        " AND " + as_of_date_clause("grids") +
                        " AND " + Grid::grid_security_clause("grids") + 
                        " AND row_locs.uuid = :uuid" +
                        " AND row_locs.version = :version" +
                        " AND row_locs.base_locale = row_locs.locale" +
                        " ORDER BY row_locs.locale", 
                         {:uuid => uuid,
                          :version => version}])
  end

  def row_begin_duplicate_exists?(row, begin_date)
    log_debug "Grid#row_begin_duplicate_exists?(begin_date=#{begin_date})"
    sql = "SELECT id" + 
          " FROM #{@db_table}" + 
          " WHERE uuid = #{quote(row.uuid)}" +
          (has_mapping? ? "" : " AND grid_uuid = #{quote(self.uuid)}") +
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
          (has_mapping? ? "" : " AND grid_uuid = #{quote(self.uuid)}") +
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
          (has_mapping? ? "" : " AND grid_uuid = #{quote(self.uuid)}")
    connection.select_value(sql).to_i
  end

  def row_title(row, full=false)
    log_debug "Grid#row_title(row=#{row}) [#{to_s}]"
    summary = (row.present? and self.has_name?) ? row.name : ""
    if not full and row.present?
      return summary if summary.length > 0
      count = 0
      load_cached_grid_structure_reference if not is_preloaded?
      @columns.each do |column|
        if [Column::REFERENCE, 
            Column::STRING, 
            Column::TEXT].include?(column.kind)
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
    if not can_select_data?
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
                          {:version => version,
                           :uuid => uuid}])
  end

  def row_select_next_version(row)
    log_debug "Grid#row_select_next_version"
    sql = "SELECT id" + 
          " FROM #{@db_table}" + 
          " WHERE uuid = #{quote(row.uuid)}" +
          " AND id != #{quote(row.id)}" +
          " AND (begin > #{quote(row.begin)}" +
          " OR (begin = #{quote(row.begin)}" +
          " AND version > #{quote(row.version)}))" +
          (has_mapping? ? "" : " AND grid_uuid = #{quote(self.uuid)}") +
          " ORDER BY begin, version"
    connection.select_value(sql).to_i
  end

  def row_select_previous_version(row)
    log_debug "Grid#row_select_previous_version"
    sql = "SELECT id" + 
          " FROM #{@db_table}" + 
          " WHERE uuid = #{quote(row.uuid)}" +
          " AND id != #{quote(row.id)}" +
          " AND (begin < #{quote(row.begin)}" +
          " OR (begin = #{quote(row.begin)}" +
          " AND version < #{quote(row.version)}))" +
          (has_mapping? ? "" : " AND grid_uuid = #{quote(self.uuid)}") +
          " ORDER BY begin DESC, version DESC"
    connection.select_value(sql).to_i
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

  def new_loc
    loc = grid_locs.new
    loc.uuid = self.uuid
    loc.version = self.version
    loc
  end
  
  def export(xml)
    log_debug "Grid#export [#{to_s}]"
    xml.grid(:title => self.name) do
      super(xml)
      xml.workspace_uuid(self.workspace_uuid, :title => self.workspace_name)
      xml.has_name(self.has_name) if self.has_name
      xml.has_description(self.has_description) if self.has_description
      Grid.all_locales(grid_locs, self.uuid, self.version).each do |loc|
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

  def copy_attributes(entity)
    log_debug "Grid#copy_attributes"
    super
    entity.workspace_uuid = self.workspace_uuid    
    entity.has_name = self.has_name    
    entity.has_description = self.has_description    
  end

  def import_attribute(xml_attribute, xml_value)
    log_debug "Grid#import_attribute(xml_attribute=#{xml_attribute}, " + 
              "xml_value=#{xml_value})"
    case xml_attribute
      when ROOT_WORKSPACE_UUID then self.workspace_uuid = xml_value
      when ROOT_HAS_NAME_UUID then self.has_name = ['true','t','1'].include?(xml_value)
      when ROOT_HAS_DESCRIPTION_UUID then self.has_description = ['true','t','1'].include?(xml_value)
    end
  end

  def import!
    log_debug "Grid#import!"
    grid = Grid.select_entity_by_uuid_version(Grid, self.uuid, self.version)
    if grid.present?
      if self.revision > grid.revision 
        log_debug "Grid#import! update"
        copy_attributes(grid)
        self.update_user_uuid = Entity.session_user_uuid
        self.updated_at = Time.now
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

  def import_loc!(loc)
    log_debug "Grid#import_loc!"
    import_loc_base!(Grid.all_locales(grid_locs, 
                                      self.uuid, 
                                      self.version), 
                                      loc)
  end

  def create_missing_loc!
    create_missing_loc_base!(Grid.all_locales(grid_locs, 
                                              self.uuid, 
                                              self.version))
  end

  def row_export(xml, row)
    log_debug "Grid#row_export(row=#{row}) [#{to_s}]"
    if row.present?
      xml.row(:title => row_title(row), :grid_uuid => self.uuid, :grid => to_s) do
        row.export(xml)
        xml.grid_uuid(self.uuid, :title => name)
        column_all.each do |column|
          xml.data(row.read_value(column), :uuid => column.uuid, :name => column.name)
        end
        row_all_locales(row.uuid, row.version).each do |loc|
          xml.locale do
            xml.base_locale(loc.base_locale)
            xml.locale(loc.locale)
            xml.name(loc.name) if loc.name.present?
            xml.description(loc.description) if loc.description.present?
          end
        end
      end
      if self.uuid == Workspace::ROOT_UUID
        grid_def = Grid::select_entity_by_uuid(Grid, Grid::ROOT_UUID)
        if grid_def.present?
          grid_def.load_cached_grid_structure
          log_debug "Grid#row_export select grids in the workspace"
          rows = grid_def.row_all([{:column_uuid => Grid::ROOT_WORKSPACE_UUID, :row_uuid => row.uuid}])
          for child in rows
            grid_def.row_export(xml, child)
          end
        end
      end
      if self.uuid == Grid::ROOT_UUID
        grid_def = Grid::select_entity_by_uuid(Grid, GridMapping::ROOT_UUID)
        if grid_def.present?
          grid_def.load_cached_grid_structure
          log_debug "Grid#row_export select mapping in the data grid"
          rows = grid_def.row_all([{:column_uuid => GridMapping::ROOT_GRID_UUID, :row_uuid => row.uuid}])
          for child in rows
            grid_def.row_export(xml, child)
          end
        end
        grid_def = Grid::select_entity_by_uuid(Grid, Column::ROOT_UUID)
        if grid_def.present?
          grid_def.load_cached_grid_structure
          log_debug "Grid#row_export select columns in the data grid"
          rows = grid_def.row_all([{:column_uuid => Column::ROOT_GRID_UUID, :row_uuid => row.uuid}])
          for child in rows
            grid_def.row_export(xml, child)
          end
        end
      end
      if self.uuid == Column::ROOT_UUID
        grid_def = Grid::select_entity_by_uuid(Grid, ColumnMapping::ROOT_UUID)
        if grid_def.present?
          grid_def.load_cached_grid_structure
          log_debug "Grid#row_export select mapping in the column"
          rows = grid_def.row_all([{:column_uuid => ColumnMapping::ROOT_COLUMN_UUID, :row_uuid => row.uuid}])
          for child in rows
            grid_def.row_export(xml, child)
          end
        end
      end
    end
  end
  
  def create_row!(row)
    log_debug "Grid#create_row!(row=#{row.inspect})"
    if not can_create_data?
      log_security_warning "Grid#create_row! Can't create data"
      return false
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
  
  def create_row_loc!(row)
    log_debug "Grid#create_row_loc!(row=#{row.inspect})"
    if not can_create_data?
      log_security_warning "Grid#create_row_loc! Can't create data"
      return false
    end
    sql = "INSERT INTO #{@db_loc_table}" +
          "(#{row_all_insert_loc_columns})" +
          " VALUES(#{row_all_insert_loc_values(row)})"
    self.id = connection.insert(sql, 
                                "#{self.class.name} Create",
                                self.class.primary_key, 
                                self.id, 
                                self.class.sequence_name)
    true
  end
  
  def update_row!(row)
    log_debug "Grid#update_row!(row=#{row.inspect})"
    if not can_update_data?
      log_security_warning "Grid#update_row! Can't update data"
      return false
    end
    row.updated_at = Time.now
    row.update_user_uuid = Entity.session_user_uuid
    sql = "UPDATE #{@db_table}" +
          " SET #{row_all_update_values(row)}" +
          " WHERE id = #{quote(row.id)}" +
          " AND lock_version = #{quote(row.lock_version)}"
    connection.update(sql, "#{self.class.name} Update")
    row.lock_version += 1
    row.make_audit(Audit::UPDATE)
    true
  end
  
  def update_row_loc!(row)
    log_debug "Grid#update_row_loc!(row=#{row.inspect})"
    if not can_update_data?
      log_security_warning "Grid#update_row_loc! Can't update data"
      return false
    end
    sql = "UPDATE #{@db_loc_table}" +
          " SET #{row_loc_update_values(row)}" +
          " WHERE id = #{quote(row.id)}" +
          " AND lock_version = #{quote(row.lock_version)}"
    connection.update(sql, "#{self.class.name} Update")
    true
  end
  
  def mapping_all
    log_debug "Grid#mapping_all"
    grid_mappings.find(:all, 
                       :conditions => 
                          [as_of_date_clause("grid_mappings")])
  end

  def mapping_select_entity_by_uuid_version(uuid, version)
    log_debug "Grid#mapping_select_entity_by_uuid_version(uuid=#{uuid}" +
              ", version=#{version})"
    grid_mappings.find(:all, 
                       :conditions => 
                          ["uuid = :uuid and version = :version", 
                          {:uuid => uuid, 
                           :version => version}])[0]
  end

  def row_initialization(row, filters)
    log_debug "Grid#row_initialization(filters=#{filters.inspect}) " + 
              "[#{to_s}]"
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
        if self.uuid == Column::ROOT_UUID and column.uuid == Column::ROOT_KIND_UUID
          row.write_value(column, Column::STRING)
        elsif self.uuid == Column::ROOT_UUID and column.uuid == Column::ROOT_DISPLAY_UUID
          row.write_value(column, 1)
        elsif self.uuid == Grid::ROOT_UUID and column.uuid == Grid::ROOT_HAS_NAME_UUID
          row.write_value(column, true)
        elsif self.uuid == Grid::ROOT_UUID and column.uuid == Grid::ROOT_HAS_DESCRIPTION_UUID
          row.write_value(column, true)
        end
      end
    end
  end

  def row_validate(row, phase)
    log_debug "Grid#row_validate(phase=#{phase}) [#{to_s}]"
    validated = true
    for column in column_all
      attribute = phase == Grid::PHASE_CREATE ?
                    column.default_physical_column :
                    column.physical_column
      value = row.read_value(column)
      log_debug "Grid#row_validate(phase=#{phase}) control" +
                " column=#{column.name}" +
                " attribute=#{attribute}<=>#{value}"
      if column.required
        log_debug "Grid#row_validate(phase=#{phase}) required"
        if value.blank?
          validated = false
          row.errors.add(attribute, I18n.t('error.required',
                                           :column => column))
        end
      end
      if column.regex.present?
        log_debug "Grid#row_validate(phase=#{phase}) regex"
        if not Regexp.new(column.regex).match(value)
          validated = false
          row.errors.add(attribute, I18n.t('error.badformat', 
                                           :column => column))
        end
      end
    end
    validated
  end

  def row_loc_validate(row, row_loc, phase)
    log_debug "Grid#row_loc_validate(phase=#{phase}) [#{to_s}]"
    validated = true
    if self.has_name?
      if row_loc.name.blank?
          validated = false
          row.errors.add(:name, I18n.t('error.mis_name'))
      end
    end
    validated
  end
  
  def load_workspace
    log_debug "Grid#load_workspace [#{to_s}]"
    self.workspace = Workspace.select_entity_by_uuid(Workspace, self.workspace_uuid)
  end
  
private

  def quote(sql)
    connection.quote(sql)
  end

  def self.all_select_columns
    "grids.id, grids.uuid, grids.version, grids.lock_version, " + 
    "grids.begin, grids.end, grids.enabled, grids.workspace_uuid, " + 
    "grids.has_name, grids.has_description, " + 
    "grids.created_at, grids.updated_at, " +
    "grids.create_user_uuid, grids.update_user_uuid, " +
    "grid_locs.base_locale, grid_locs.locale, " +
    "grid_locs.name, grid_locs.description"
  end
  
  def column_all_select_columns
    "columns.id, columns.uuid, columns.version, columns.lock_version, " +
    "columns.begin, columns.end, columns.enabled, " +
    "columns.display, columns.kind, " + 
    "columns.grid_reference_uuid, " + 
    "columns.required, columns.regex, " + 
    "columns.created_at, columns.updated_at, " + 
    "columns.create_user_uuid, columns.update_user_uuid, " + 
    "column_locs.base_locale, column_locs.locale, " +
    "column_locs.name, column_locs.description"
  end
  
  def row_all_select_columns
    columns = ""
    column_all.each do |column|
      columns << ", rows." + column.physical_column
    end
    "#{quote(self.uuid)} grid_uuid" +
    ", rows.id, rows.uuid, rows.version, rows.lock_version" + 
    ", rows.begin, rows.end, rows.enabled" + 
    ", rows.created_at, rows.updated_at" +
    ", rows.create_user_uuid, rows.update_user_uuid" +
    (has_translation? ? 
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
    column_all.each do |column|
      columns << ", " + column.physical_column
    end
    "uuid, version, lock_version, begin, end, enabled" + 
    (has_mapping? ? "" : ", grid_uuid") + 
    ", created_at, updated_at, create_user_uuid, update_user_uuid" + 
    columns
  end
  
  def row_all_insert_loc_columns
    "uuid, version, lock_version, locale, base_locale, name, description" 
  end
  
  def row_all_insert_values(row)
    columns = ""
    column_all.each do |column|
      columns << ", #{quote(row.read_value(column))}"
    end
    "#{quote(row.uuid)}" +
    ", #{quote(row.version)}" +
    ", #{quote(row.lock_version)}" +
    ", #{quote(row.begin)}" +
    ", #{quote(row.end)}" +
    ", #{quote(row.enabled)}" +
    (has_mapping? ? "" : ", #{quote(uuid)}") + 
    ", #{quote(row.created_at)}" +
    ", #{quote(row.updated_at)}" +
    ", #{quote(row.create_user_uuid)}" +
    ", #{quote(row.update_user_uuid)}" +
    columns
  end
  
  def row_all_insert_loc_values(row)
    "#{quote(row.uuid)}" +
    ", #{quote(row.version)}" +
    ", #{quote(row.lock_version)}" +
    ", #{quote(row.locale)}" +
    ", #{quote(row.base_locale)}" +
    ", #{quote(row.name)}" +
    ", #{quote(row.description)}"
  end
  
  def row_all_update_values(row)
    columns = ""
    column_all.each do |column|
      columns << ", #{column.physical_column}=#{quote(row.read_value(column))}"
    end
    "lock_version=#{quote(row.lock_version+1)}" +
    ", begin=#{quote(row.begin)}" +
    ", end=#{quote(row.end)}" +
    ", enabled=#{quote(row.enabled)}" +
    ", updated_at=#{quote(row.updated_at)}" +
    ", update_user_uuid=#{quote(row.update_user_uuid)}" +
    columns
  end
  
  def row_loc_update_values(row)
    "lock_version=#{quote(row.lock_version+1)}" +
    ", base_locale=#{quote(row.base_locale)}" +
    ", name=#{quote(row.name)}" +
    ", description=#{quote(row.description)}"
  end
  
  def grid_mapping_read_select_columns
    "grid_mappings.uuid, grid_mappings.version, " +
    "grid_mappings.begin, grid_mappings.end, " +
    "grid_mappings.db_table, grid_mappings.db_loc_table"
  end
  
  def grid_mapping_read
    grid_mappings.find(:first, 
                       :select => grid_mapping_read_select_columns,
                       :conditions => 
                              [as_of_date_clause("grid_mappings")])
  end

  # Loads in memory information about database mapping
  def load_cached_mapping
    log_debug "Grid#load_cached_mapping [#{to_s}]"
    grid_mapping = grid_mapping_read
    if grid_mapping.present?
      @db_table = grid_mapping.db_table if grid_mapping.db_table.present?
      @db_loc_table = grid_mapping.db_loc_table if grid_mapping.db_loc_table.present?
    else
      @db_table = "rows"
      @db_loc_table = "row_locs"
    end
    @has_mapping = grid_mapping.present?
  end

  # Loads in memory information about columns
  def load_columns(filters, skip_reference=false, skip_mapping=false)
    log_debug "Grid#load_columns " +
              "filters=#{filters},"+
              "skip_reference=#{skip_reference}," +
              "skip_mapping=#{skip_mapping}) [#{to_s}]"
    @all_columns = columns.find(:all,
                 :joins => :column_locs,
                 :select => column_all_select_columns,
                 :conditions =>
                      [as_of_date_clause("columns") +
                       " AND column_locs.version = columns.version" +
                       " AND " + locale_clause("column_locs")],
                 :order => "columns.id")
    @columns = Array.new(@all_columns)
    index_reference = index_date = index_integer = index_decimal = index_string = 0
    column_all.each do |column|
      unless column.is_preloaded?
        log_debug "Grid#load_columns column=#{column.name}"
        case column.kind
          when Column::REFERENCE then
            index_reference += 1 
            number = index_reference
          when Column::DATE then 
            index_date += 1 
            number = index_date
          when Column::INTEGER then 
            index_integer += 1 
            number = index_integer
          when Column::DECIMAL then 
            index_decimal += 1 
            number = index_decimal
          else 
            index_string += 1 
            number = index_string
        end
        log_debug "Grid#load_columns number=#{number}"
        column.load_cached_information(self,
                                       number,
                                       skip_reference,
                                       skip_mapping)
        log_debug "Grid#load_columns column.physical_column=#{column.physical_column}"
        @columns.delete(column) if column.display.blank? or column.display == 0
        unless filters.nil?
          filters.each{|filter| @columns.delete(column) if filter[:column_uuid] == column.uuid}
        end
      end
    end
    @columns.sort! {|a, b| a.display <=> b.display}
  end
  
  def load_security_workspace(uuid)
    log_debug "Grid#load_security_workspace [#{to_s}]"
    
    sql = "SELECT workspace_sharings.role_uuid" + 
          " FROM workspace_sharings" + 
          " WHERE workspace_sharings.workspace_uuid = '#{uuid}'" + 
          " AND workspace_sharings.user_uuid = '#{Entity.session_user_uuid}'" +
          " AND " + as_of_date_clause("workspace_sharings") +
          " LIMIT 1"

    security = Grid.find_by_sql([sql])[0]
    log_debug "Grid#load_security_workspace (1) security=#{security.inspect}"
    return security.role_uuid if not security.nil?

    sql = "SELECT default_role_uuid, create_user_uuid" + 
          " FROM workspaces workspace_security" + 
          " WHERE workspace_security.uuid = '#{uuid}'" + 
          " AND (" +
          "  workspace_security.public = 't'" +
          "  OR workspace_security.default_role_uuid in ('#{Role::ROLE_READ_ONLY_UUID}', '#{Role::ROLE_READ_WRITE_UUID}', '#{Role::ROLE_READ_WRITE_ALL_UUID}', '#{Role::ROLE_TOTAL_CONTROL_UUID}')" +
          "  OR workspace_security.create_user_uuid = '#{Entity.session_user_uuid}'" +
          " )" +
          " AND " + as_of_date_clause("workspace_security") +
          " LIMIT 1"
    security = Grid.find_by_sql([sql])[0]
    log_debug "Grid#load_security_workspace (2) security=#{security.inspect}"
    return Role::ROLE_TOTAL_CONTROL_UUID if not security.nil? and 
                      security.create_user_uuid == Entity.session_user_uuid 
    return security.default_role_uuid if not security.nil?

    nil
  end

  def load_security
    log_debug "Grid#load_security [#{to_s}]"
    sql = "SELECT workspace_sharings.role_uuid" + 
          " FROM workspace_sharings" + 
          " WHERE workspace_sharings.workspace_uuid = '#{self.workspace_uuid}'" + 
          " AND workspace_sharings.user_uuid = '#{Entity.session_user_uuid}'" +
          " AND " + as_of_date_clause("workspace_sharings") +
          " LIMIT 1"
    security = Grid.find_by_sql([sql])[0]
    if security.present?
      @can_select_data = [Role::ROLE_READ_ONLY_UUID, Role::ROLE_READ_WRITE_UUID, Role::ROLE_READ_WRITE_ALL_UUID, Role::ROLE_TOTAL_CONTROL_UUID].include?(security.role_uuid)
      @can_create_data = [Role::ROLE_READ_WRITE_UUID, Role::ROLE_READ_WRITE_ALL_UUID, Role::ROLE_TOTAL_CONTROL_UUID].include?(security.role_uuid)
      @can_update_data = [Role::ROLE_READ_WRITE_UUID, Role::ROLE_READ_WRITE_ALL_UUID, Role::ROLE_TOTAL_CONTROL_UUID].include?(security.role_uuid)
    else
      sql = "SELECT default_role_uuid, public" + 
            " FROM workspaces workspace_security" + 
            " WHERE workspace_security.uuid = '#{self.workspace_uuid}'" + 
            " AND (" +
            "  workspace_security.public = 't'" +
            "  OR workspace_security.default_role_uuid in ('#{Role::ROLE_READ_ONLY_UUID}', '#{Role::ROLE_READ_WRITE_UUID}', '#{Role::ROLE_READ_WRITE_ALL_UUID}', '#{Role::ROLE_TOTAL_CONTROL_UUID}')" +
            "  OR workspace_security.create_user_uuid = '#{Entity.session_user_uuid}'" +
            " )" +
            " AND " + as_of_date_clause("workspace_security") +
            " LIMIT 1"
      security = Grid.find_by_sql([sql])[0]
      if security.present? and security.default_role_uuid.present?
        @can_select_data = [Role::ROLE_READ_ONLY_UUID, Role::ROLE_READ_WRITE_UUID, Role::ROLE_READ_WRITE_ALL_UUID, Role::ROLE_TOTAL_CONTROL_UUID].include?(security.default_role_uuid)
        @can_create_data = [Role::ROLE_READ_WRITE_UUID, Role::ROLE_READ_WRITE_ALL_UUID, Role::ROLE_TOTAL_CONTROL_UUID].include?(security.default_role_uuid)
        @can_update_data = [Role::ROLE_READ_WRITE_UUID, Role::ROLE_READ_WRITE_ALL_UUID, Role::ROLE_TOTAL_CONTROL_UUID].include?(security.default_role_uuid)
      elsif security.present? and security.public
        @can_select_data = true
        @can_create_data = false
        @can_update_data = false
      end
    end
    log_debug "Grid#load_security " +
              "@can_select_data=#{@can_select_data}," +
              "@can_create_data=#{@can_create_data}," +
              "@can_update_data=#{@can_update_data}," +
              "[#{to_s}]"
  end
end

class GridLoc < EntityLoc
  validates_presence_of :name
end