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
  ROOT_DISPLAY_UUID = '4c766f00-06cd-012d-c1fb-0026b0d63708'
  ROOT_REFERENCE_UUID = '5a5107e0-0990-012d-e81a-4417fe7fde95'
  ROOT_DATA_KIND_UUID = '5a2e26e0-ea31-012c-1074-00166f92f624'
  ROOT_REQUIRED_UUID = '8df46b61-1da3-012d-46db-4417fe7fde95'
  ROOT_REGEX_UUID = '68792a30-1da5-012d-2556-4417fe7fde95'

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
  validates_presence_of :grid_uuid, :kind
  validates_associated :grid
  attr_reader :physical_column,
              :default_physical_column,
              :grid_reference,
              :workspace_reference
  
  def is_preloaded?
    physical_column.present?
  end
  
  def load_cached_information(grid, number, skip_reference, skip_mapping)
    log_debug "Column#load_cached_information(grid=#{grid.to_s})"
    case self.kind
      when REFERENCE then 
        @default_physical_column = "row_uuid" + number.to_s
      when DATE then 
        @default_physical_column = "date" + number.to_s
      when INTEGER then 
        @default_physical_column = "integer" + number.to_s
      when DECIMAL then 
        @default_physical_column = "float" + number.to_s
      else 
        @default_physical_column = "value" + number.to_s
    end
    db_column = column_mapping_column
    if not skip_mapping and db_column.present?
      @physical_column = db_column
    else
      @physical_column = default_physical_column
    end
    log_debug "Column#load_cached_information " + 
              "physical_column=#{physical_column}"
    if not skip_reference
      if self.kind == REFERENCE and 
         self.grid_reference_uuid.present?
        # this is used to avoid circular references
        log_debug "Column#load_cached_information reference " +
                  "grid_reference_uuid=#{self.grid_reference_uuid}"
        if self.grid_reference_uuid == grid.uuid
          log_error "Column#load_cached_information circular reference " +
                    "for data grid '#{self.grid_reference_uuid}'"
          @grid_reference = grid
          @workspace_reference = 
            Workspace.select_entity_by_uuid(Workspace,
                                            grid.workspace_uuid) if grid.present?
        else
          @grid_reference = 
            Grid.select_entity_by_uuid(Grid, self.grid_reference_uuid)
          log_debug "Column#load_cached_information " +
                    "@grid_reference=#{@grid_reference.to_s}"
          if @grid_reference.present? 
            @grid_reference.load_cached_grid_structure_reference if not @grid_reference.is_preloaded?
            @workspace_reference = Workspace.select_entity_by_uuid(Workspace,
                                                                   @grid_reference.workspace_uuid)
          end
        end
      end
    end
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
  
  def import_attribute(xml_attribute, xml_value)
    log_debug "Column#import_attribute(xml_attribute=#{xml_attribute}, " + 
              "xml_value=#{xml_value})"
    case xml_attribute
      when ROOT_GRID_UUID then self.grid_uuid = xml_value
      when ROOT_REFERENCE_UUID then self.grid_reference_uuid = xml_value
      when ROOT_KIND_UUID then self.kind = xml_value
      when ROOT_DISPLAY_UUID then self.display = xml_value.to_i
      when ROOT_REQUIRED_UUID then self.required = ['true','t','1'].include?(xml_value)
      when ROOT_REGEX_UUID then self.regex = xml_value
    end
  end
  
  def copy_attributes(entity)
    log_debug "Column#copy_attributes"
    super
    entity.grid_uuid = self.grid_uuid
    entity.grid_reference_uuid = self.grid_reference_uuid
    entity.kind = self.kind
    entity.display = self.display
    entity.required = self.required
    entity.regex = self.regex
  end

  def import!
    log_debug "Column#import!"
    log_error "Can't import column when there is no grid reference" if grid.nil?
    column = grid.column_select_entity_by_uuid_version(self.uuid, self.version)
    if column.present?
      if self.revision > column.revision 
        log_debug "Column#import! update"
        copy_attributes(column)
        make_audit(Audit::IMPORT)
        column.save!
        column.update_dates!(grid.columns, update_user_uuid)
        return "updated"
      else
        log_debug "Column#import! skip update"
        return "skipped"
      end
    else
      log_debug "Column#import! new"
      make_audit(Audit::IMPORT)
      save!
      update_dates!(grid.columns)
      return "inserted"
    end
    ""
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