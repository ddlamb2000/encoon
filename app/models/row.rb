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
class Row < Entity
  XML_TAG = 'row'
  
  belongs_to :grid, :foreign_key => "grid_uuid", :primary_key => "uuid"
  has_many :row_locs, :foreign_key => "uuid", :primary_key => "uuid"
  has_many :row_attachments, :foreign_key => "uuid", :primary_key => "uuid"
  has_many :row_passwords, :foreign_key => "uuid", :primary_key => "uuid"
  validates_presence_of :grid_uuid
  validates_associated :grid
  
  @initialization = false
  
  def initialization
    defaults
    @initialization = true
  end

  def initialization?
    @initialization
  end

  def before_destroy
    log_debug "Row#before_destroy"
    super
    RowLoc.destroy_all(["uuid = :uuid and version = :version", 
                       {:uuid => self.uuid, :version => self.version}])
    row_attachments.destroy_all
    row_passwords.destroy_all
  end
  
  def to_s
    name
  end

  def name
    attribute_present?(:name) ? 
      read_attribute(:name) : 
      ""
  end

  def title
    attribute_present?(:name) ? 
      read_attribute(:name) : 
      (grid.present? ? 
        grid.row_title(self) : 
        "")
  end

  def description
    attribute_present?(:description) ? 
      read_attribute(:description) : 
      ""
  end

  def read_value(column)
    read_attribute(@initialization ? 
      column.default_physical_column : 
      column.physical_column)
  end
  
  def write_value(column, value)
    log_debug "Row#write_value column.physical_column=#{column.name}, " +
              "value=#{value.inspect}"
    send((@initialization ? 
           column.default_physical_column : 
           column.physical_column) + 
         '=', 
         value)
  end
  
  def read_referenced_name(column)
    value = read_value(column)
    if value.present? 
      if column.kind == Column::REFERENCE and 
         column.grid_reference_uuid.present?
        log_debug "Row#read_referenced_name value=#{value}, " +
                  "column.grid_reference_uuid=#{column.grid_reference_uuid}"
        grid = column.loaded_grid_reference
        if grid.present?
          grid.load_cached_grid_structure
          return grid.select_reference_row_name(value)  
        end
      end
      value
    else
      ""
    end
  end

  def read_referenced_description(column)
    if column.kind == Column::REFERENCE and 
       column.grid_reference_uuid.present?
      value = read_value(column)
      if value.present? 
        log_debug "Row#read_referenced_description value=#{value}, " +
                  "column.grid_reference_uuid=#{column.grid_reference_uuid}"
        grid = column.loaded_grid_reference
        if grid.present?
          grid.load_cached_grid_structure
          return grid.select_reference_row_description(value)  
        end
      end
    end
    ""
  end

  def self.select_grid_cast(grid_uuid, row_uuid)
    log_debug "Row#select_grid_cast grid_uuid=#{grid_uuid}, " + 
              "row_uuid=#{row_uuid}"
    if grid_uuid == Grid::ROOT_UUID
      grid = Grid.select_entity_by_uuid(Grid, row_uuid)
      grid.load_cached_grid_structure if grid.present?
      grid
    end
  end

  def self.loc_select_columns
    "row_locs.id, row_locs.uuid, " + 
    "row_locs.version, row_locs.lock_version, " + 
    "row_locs.base_locale, row_locs.locale, " +
    "row_locs.name, row_locs.description"
  end
  
  def self.all_locales(collection, uuid, version)
    collection.find(:all, 
                    :select => self.loc_select_columns,
                    :conditions => 
                           ["row_locs.uuid = :uuid " +
                            "AND row_locs.version = :version " +
                            "AND row_locs.locale = row_locs.base_locale", 
                           {:uuid => uuid, 
                            :version => version}], 
                    :order => "row_locs.locale")
  end

  def new_loc
    loc = row_locs.new
    loc.uuid = self.uuid
    loc.version = self.version
    loc
  end

  def import(xml_attribute, xml_value)
    log_debug "Row#import(xml_attribute=#{xml_attribute}, " + 
               "xml_value=#{xml_value})"
    case xml_attribute
      when 'grid_uuid' then self.grid_uuid = xml_value
      else super
    end
  end

  def import_attribute(xml_attribute, xml_value)
    log_debug "Row#import_attribute(xml_attribute=#{xml_attribute}, " + 
              "xml_value=#{xml_value})"
    self.initialization
    grid.load_cached_grid_structure_reference if not grid.is_preloaded? 
    grid.column_all.each do |column|
      log_debug "Row#import_attribute column=#{column}, " + 
                "column.uuid=#{column.uuid}"
      if xml_attribute == column.uuid
        write_value(column, xml_value)
      end
    end
  end

  def copy_attributes(entity)
    log_debug "Row#copy_attributes"
    super
    entity.grid_uuid = self.grid_uuid
    grid.load_cached_grid_structure_reference if not grid.is_preloaded? 
    grid.column_all.each do |column|
      log_debug "Row#copy_attributes column=#{column}"
      write_value(column, self.read_value(column))
    end
  end

  def import!
    log_debug "Row#import!"
    raise "Can't import row when there is no grid reference" if grid.nil?
    row = grid.row_select_entity_by_uuid_version(self.uuid, self.version)
    if row.present?
      if self.revision > row.revision 
        log_debug "Row#import! update"
        copy_attributes(row)
        self.update_user_uuid = Grid.session_user_uuid    
        make_audit(Audit::IMPORT)
        grid.update_row!(self)
        grid.row_update_dates!(self.uuid)
        return true
      else
        log_debug "Row#import! skip update"
      end
    else
      log_debug "Row#import! new"
      self.create_user_uuid = self.update_user_uuid = Grid.session_user_uuid    
      self.created_at = self.updated_at = Time.now    
      make_audit(Audit::IMPORT)
      grid.create_row!(self)
      grid.row_update_dates!(self.uuid)
      return true
    end
    false
  end

  def import_loc!(loc)
    log_debug "Row#import_loc!(loc=#{loc})"
    raise "Can't import row loc when there is no grid reference" if grid.nil?
    updated = 0
    grid.row_loc_select_entity_by_uuid(self.uuid, 
                                       self.version).each do |row_loc|
      if row_loc.base_locale == loc.base_locale
        log_debug "Row#import_loc! update"
        loc.copy_attributes(row_loc)
        grid.update_row_loc!(row_loc)
        updated += 1
      end
    end
    if updated == 0
      log_debug "Row#import_loc! new"
      grid.create_row_loc!(loc)
    end
  end

  def create_missing_loc!
    raise "Can't create row loc when there is no grid reference" if grid.nil?
    base_locs = []
    base_loc = nil
    foundI18n = false
    grid.row_loc_select_entity_by_uuid(self.uuid, self.version).each do |loc|
      base_locs << loc.base_locale
      if not foundI18n
        base_loc = loc
        foundI18n = (loc.base_locale == I18n.locale.to_s)
      end
    end
    if base_loc.present?
      LANGUAGES.each do |lang, locale|
        if (base_locs.find {|value| locale.to_s == value}).nil?
          loc = new_loc        
          base_loc.copy_attributes(loc)
          loc.locale = locale.to_s
          loc.base_locale = base_loc.base_locale
          grid.create_row_loc!(loc)
          log_debug "Row#create_missing! new, locale=#{locale.to_s}"
        end
      end
    end
  end
  
  def has_document?
    for attachment in row_attachments
      return true if attachment.document?
    end
    false
  end

  def has_photo?
    for attachment in row_attachments
      return true if attachment.photo?
    end
    false
  end

  def first_photo
    for attachment in row_attachments
      return attachment if attachment.photo?
    end
    nil
  end

  def has_attachment?
    for attachment in row_attachments
      return true if attachment.document? or attachment.photo?
    end
    false
  end
  
  def has_password?
    for password in row_passwords
      return true if password.salt.present? and password.password.present? 
    end
    false
  end

  def remove_attachment!(input_file)
    for attachment in row_attachments
      if attachment.content_type == input_file.content_type.chomp and 
         attachment.file_name == input_file.original_filename
        attachment.destroy
      end
    end
  end

  def remove_password!
    for password in row_passwords
      password.delete
    end
  end
end

class RowLoc < EntityLoc
end