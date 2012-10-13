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
  belongs_to :grid, :foreign_key => "grid_uuid", :primary_key => "uuid"
  belongs_to :grid_reference, :class_name => "Grid", :foreign_key => "grid_reference_uuid", :primary_key => "uuid"
  has_many :column_locs, :foreign_key => "uuid", :primary_key => "uuid"
  has_many :column_mappings, :foreign_key => "column_uuid", :primary_key => "uuid"
  validates_presence_of :grid_uuid, :kind
  validates_associated :grid
  attr_reader :physical_column, :default_physical_column, :grid_reference, :workspace_reference
  
  def loaded ; physical_column.present? ; end
  
  def load_cached_information(grid, number, skip_reference, skip_mapping, db_column)
    case self.kind
      when COLUMN_TYPE_REFERENCE then 
        @default_physical_column = "row_uuid" + number.to_s
      when COLUMN_TYPE_DATE then 
        @default_physical_column = "date" + number.to_s
      when COLUMN_TYPE_INTEGER then 
        @default_physical_column = "integer" + number.to_s
      when COLUMN_TYPE_DECIMAL then 
        @default_physical_column = "float" + number.to_s
      else 
        @default_physical_column = "value" + number.to_s
    end
    @physical_column = (skip_mapping or db_column.nil?) ? @default_physical_column : db_column
    log_debug "Column#load_cached_information grid=#{grid.to_s}, name=#{self.name}, physical_column=#{physical_column}"
    if not(skip_reference)
      if self.kind == COLUMN_TYPE_REFERENCE and 
         self.grid_reference_uuid.present?
        # this is used to avoid circular references
        log_debug "Column#load_cached_information reference " +
                  "grid_reference_uuid=#{self.grid_reference_uuid}"
        if self.grid_reference_uuid == grid.uuid
          log_error "Column#load_cached_information circular reference " +
                    "for data grid '#{self.grid_reference_uuid}'"
          @grid_reference = grid
          @workspace_reference = Workspace.select_entity_by_uuid(Workspace,
                                                                 grid.workspace_uuid) if grid.present?
        else
          @grid_reference = 
            Grid.select_entity_by_uuid(Grid, self.grid_reference_uuid)
          log_debug "Column#load_cached_information @grid_reference=#{@grid_reference.to_s}"
          if @grid_reference.present? 
            @grid_reference.load if not(@grid_reference.loaded)
            @workspace_reference = Workspace.select_entity_by_uuid(Workspace,
                                                                   @grid_reference.workspace_uuid)
          end
        end
      end
    end
  end
  
  def new_loc
    loc = column_locs.new
    loc.uuid = self.uuid
    loc.version = self.version
    loc
  end
  
  def import_attribute(xml_attribute, xml_value)
    log_debug "Column#import_attribute(#{xml_attribute}, #{xml_value})"
    case xml_attribute
      when COLUMN_GRID_UUID then self.grid_uuid = xml_value
      when COLUMN_REFERENCE_UUID then self.grid_reference_uuid = xml_value
      when COLUMN_KIND_UUID then self.kind = xml_value
      when COLUMN_DISPLAY_UUID then self.display = xml_value.to_i
      when COLUMN_REQUIRED_UUID then self.required = ['true','t','1'].include?(xml_value)
      when COLUMN_REGEX_UUID then self.regex = xml_value
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
        self.update_user_uuid = Entity.session_user_uuid
        self.updated_at = Time.now
        make_audit(Audit::IMPORT)
        column.save!
        column.update_dates!(grid.columns)
        return "updated"
      else
        log_debug "Column#import! skip update"
        return "skipped"
      end
    else
      log_debug "Column#import! new"
      self.create_user_uuid = self.update_user_uuid = Entity.session_user_uuid
      self.created_at = self.updated_at = Time.now
      make_audit(Audit::IMPORT)
      save!
      update_dates!(grid.columns)
      return "inserted"
    end
    ""
  end

  def import_loc!(loc)
    log_debug "Column#import_loc!(loc=#{loc})"
    import_loc_base!(Column.locales(column_locs, self.uuid, self.version), loc)
  end

  # Creates local row for all the installed languages
  # that is not created yet for the given collection.
  # This insures on row exists for any installed language.
  def create_missing_loc!
    log_debug "Column#create_missing_loc!"
    create_missing_loc_base!(Column.locales(column_locs, self.uuid, self.version))
  end

  def column_mapping_select_entity_by_uuid_version(uuid, version)
    log_debug "Column#column_mapping_select_entity_by_uuid_version(" + 
              "uuid=#{uuid}, version=#{version})"
    column_mappings.find(:all, 
                         :conditions => 
                          ["uuid = :uuid and version = :version", 
                          {:uuid => uuid, :version => version}])[0]
  end

private

  def self.loc_select_columns
    "column_locs.id, column_locs.uuid, " + 
    "column_locs.version, column_locs.lock_version, " + 
    "column_locs.base_locale, column_locs.locale, " +
    "column_locs.name, column_locs.description"
  end
end