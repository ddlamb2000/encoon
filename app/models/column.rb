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
class Column < Entity
  ROOT_UUID = '03bfafe0-ea31-012c-105a-00166f92f624'
  ROOT_GRID_UUID = '461e5940-ea31-012c-106b-00166f92f624'
  ROOT_KIND_UUID = '50270360-ea31-012c-106f-00166f92f624'
  ROOT_NUMBER_UUID = '42e7f000-06cd-012d-c1f8-0026b0d63708'
  ROOT_DISPLAY_UUID = '4c766f00-06cd-012d-c1fb-0026b0d63708'
  ROOT_REFERENCE_UUID = '5a5107e0-0990-012d-e81a-4417fe7fde95'
  
  ROOT_DATA_KIND_UUID = '5a2e26e0-ea31-012c-1074-00166f92f624'
  ROOT_GRID_DISPLAY_OPTION_UUID = 'f4b48831-1df3-012d-8b41-4417fe7fde95'
  
  XML_TAG = 'column'

  STRING = "603286d0-ea31-012c-1079-00166f92f624"
  TEXT = "688c4870-ea31-012c-107d-00166f92f624"
  DATE = "6d2c8470-ea31-012c-1081-00166f92f624"
  INTEGER = "f33188f0-06cc-012d-c1db-0026b0d63708"
  DECIMAL = "f766af30-06cc-012d-c1df-0026b0d63708"
  BOOLEAN = "ff2dd910-06cc-012d-c1e3-0026b0d63708"
  REFERENCE = "04950150-06cd-012d-c1e7-0026b0d63708"
  HYPERLINK = "0a0a5050-06cd-012d-c1eb-0026b0d63708"
  PASSWORD = "91add730-0d0b-012d-4ae1-4417fe7fde95"
  
  belongs_to :grid, :foreign_key => "grid_uuid", :primary_key => "uuid"
  belongs_to :grid_reference, :class_name => "Grid", :foreign_key => "grid_reference_uuid", :primary_key => "uuid"
  has_many :column_locs, :foreign_key => "uuid", :primary_key => "uuid"
  has_many :column_mappings, :foreign_key => "column_uuid", :primary_key => "uuid"
  validates_presence_of :grid_uuid, :number, :display, :kind
  validates :number, :inclusion => { :in => 1..20 }
  validates_associated :grid
  
  def before_destroy
    log_debug "Column#before_destroy [column #{to_s}]"
    super
    ColumnLoc.destroy_all(["uuid = :uuid AND version = :version", 
                          {:uuid => self.uuid, :version => self.version}])
    if grid.column_all_versions(self.uuid).length == 0
      log_debug "Column#before_destroy remove_orphans"
      column_mappings.destroy_all
    end
  end
  
  def is_preloaded?
    @physical_column.present?
  end
  
  def has_mapping?
    @has_mapping
  end
  
  def load_cached_information(grid_uuid, grid, skip_reference, skip_mapping)
    case self.kind
      when REFERENCE then 
        @default_physical_column = "row_uuid" + self.number.to_s
      when DATE then 
        @default_physical_column = "date" + self.number.to_s
      when INTEGER then 
        @default_physical_column = "integer" + self.number.to_s
      when DECIMAL then 
        @default_physical_column = "float" + self.number.to_s
      else 
        @default_physical_column = "value" + self.number.to_s
    end
    db_column = column_mapping_column
    if not skip_mapping and db_column.present?
      @physical_column = db_column
    else
      @physical_column = @default_physical_column
    end
    @has_mapping = db_column.present?
    if not skip_reference
      if self.kind == REFERENCE and 
         self.grid_reference_uuid.present?
        # this is used to avoid circular references
        if self.grid_reference_uuid == grid_uuid
          log_error "Column#load_cached_information circular reference " +
                    "for data grid '#{self.grid_reference_uuid}'"
          @loaded_grid_reference = grid
        else
          @loaded_grid_reference = 
            Grid.select_entity_by_uuid(Grid, self.grid_reference_uuid)
          if @loaded_grid_reference.present? and 
              not @loaded_grid_reference.is_preloaded?
            @loaded_grid_reference.load_cached_grid_structure_reference
          end
        end
      end
    end
  end
  
  def physical_column
    if @physical_column.nil?
      log_error "Column #{to_s} for data grid #{grid.to_s} isn't preloaded."
    end
    @physical_column 
  end
  
  def default_physical_column
    if @default_physical_column.nil?
      log_error "Column #{to_s} for data grid #{grid.to_s} isn't preloaded."
    end
    @default_physical_column 
  end
  
  def loaded_grid_reference
    if @loaded_grid_reference.nil?
      log_error "Column#loaded_grid_reference column #{to_s} isn't preloaded."
    end
    @loaded_grid_reference
  end
  
  # Returns the name of the grid used as a reference
  def grid_reference_name
    grid = Grid.select_entity_by_uuid(Grid, grid_reference_uuid)
    grid.present? ? grid.reference_name : ""
  end

  def self.all_locales(collection, uuid, version)
    collection.find(:all, 
                    :select => self.loc_select_columns,
                    :conditions => 
                           ["column_locs.uuid = :uuid " +
                            "AND column_locs.version = :version " +
                            "AND column_locs.locale = column_locs.base_locale", 
                           {:uuid => uuid, 
                            :version => version}], 
                    :order => "column_locs.locale")
  end

  def new_loc
    loc = column_locs.new
    loc.uuid = self.uuid
    loc.version = self.version
    loc
  end
  
  def import(xml_attribute, xml_value)
    log_debug "Column#import(xml_attribute=#{xml_attribute}, " + 
              "xml_value=#{xml_value})"
    case xml_attribute
      when 'grid_uuid' then self.grid_uuid = xml_value
      when 'grid_reference_uuid' then
        self.grid_reference_uuid = xml_value
      when 'kind' then self.kind = xml_value
      when 'number' then self.number = xml_value.to_i
      when 'display' then self.display = xml_value.to_i
      when 'required' then self.required = (xml_value == 'true')
      when 'length' then self.length = xml_value.to_i
      when 'decimals' then self.decimals = xml_value.to_i
      when 'regex' then self.regex = xml_value
      else super
    end
  end
  
  def copy_attributes(entity)
    log_debug "Column#copy_attributes"
    super
    entity.grid_uuid = self.grid_uuid
    entity.grid_reference_uuid = self.grid_reference_uuid
    entity.kind = self.kind
    entity.number = self.number
    entity.display = self.display
    entity.required = self.required
    entity.length = self.length
    entity.decimals = self.decimals
    entity.regex = self.regex
  end

  def import!
    log_debug "Column#import!"
    log_error "Can't import column when there is no grid reference" if grid.nil?
    column = grid.column_select_entity_by_uuid_version(self.uuid)
    if column.present?
      if self.revision > column.revision 
        logger.debug "Column#import! update"
        copy_attributes(column)
        self.update_user_uuid = Entity.session_user_uuid    
        make_audit(Audit::IMPORT)
        column.save!
        column.update_dates!(grid.columns, update_user_uuid)
        return true
      else
        log_debug "Column#import! skip update"
      end
    else
      log_debug "Column#import! new"
      self.create_user_uuid = self.update_user_uuid = Entity.session_user_uuid    
      self.created_at = self.updated_at = Time.now    
      make_audit(Audit::IMPORT)
      save!
      update_dates!(grid.columns)
      return true
    end
    false
  end

  def import_loc!(loc)
    log_debug "Column#import_loc!(loc=#{loc})"
    import_loc_base!(
      Column.all_locales(column_locs, self.uuid, self.version), loc)
  end

  def create_missing_loc!
    create_missing_loc_base!(
      Column.all_locales(column_locs, self.uuid, self.version))
  end

  def column_mapping_all
    log_debug "Column#column_mapping_all"
    column_mappings.find(:all, 
                         :conditions => 
                            [as_of_date_clause("column_mappings")])
  end

  def column_mapping_select_entity_by_uuid_version(uuid, version)
    log_debug "Column#column_mapping_select_entity_by_uuid_version(" + 
              "uuid=#{uuid}, version=#{version})"
    column_mappings.find(:all, 
                         :conditions => 
                          ["uuid = :uuid and version = :version", 
                          {:uuid => uuid, 
                           :version => version}])[0]
  end

private

  def column_mapping_read_select_columns
    "column_mappings.uuid, column_mappings.version, " +
    "column_mappings.begin, column_mappings.end, " +
    "column_mappings.db_column"
  end
  
  def column_mapping_read
    log_debug "Column#column_mapping_read [column #{to_s}]"
    column_mappings.find(:first, 
                         :select => column_mapping_read_select_columns,
                         :conditions => 
                              [Grid.as_of_date_clause("column_mappings")])
  end

  def column_mapping_column
    column_mapping = column_mapping_read
    column_mapping.present? ? column_mapping.db_column : nil
  end

  def self.loc_select_columns
    "column_locs.id, column_locs.uuid, " + 
    "column_locs.version, column_locs.lock_version, " + 
    "column_locs.base_locale, column_locs.locale, " +
    "column_locs.name, column_locs.description"
  end
end

class ColumnLoc < EntityLoc
  validates_presence_of :name
end