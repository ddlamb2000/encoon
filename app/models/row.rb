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
  belongs_to :grid, :foreign_key => "grid_uuid", :primary_key => "uuid"
  has_many :row_locs, :foreign_key => "uuid", :primary_key => "uuid"
  has_many :attachments, :foreign_key => "uuid", :primary_key => "uuid", :order => "document_file_name"

  # Basic controls.
  validates_presence_of :grid_uuid
  validates_associated :grid

  attr_reader :initialized
  
  @initialized = false
  
  def initialization
    defaults
    @initialized = true
  end

  # Returns the title information associated to the row if exists.
  # If a name column doesn't exist, the title is set by the grid the row is attached to.
  def title
    attribute_present?(:name) ? read_attribute(:name) : (grid.present? ? grid.row_title(self) : "")
  end

  # Reads the value of a given row column.
  # Reads the value from a generic column
  # or from the physical column of a table when the category is mapped to a table. 
  def read_value(column)
    read_attribute(@initialized ? column.default_physical_column : column.physical_column)
  end

  # Writes the given value in the data row column.
  # Writes the value in a generic column that depends on the type of value,
  # or in the physical column of a table when the category is mapped to a table. 
  def write_value(column, value)
    output = @initialized ? column.default_physical_column : column.physical_column
    if column.kind == COLUMN_TYPE_DATE and value.is_a?(String)
      begin
        log_debug "Row#write_value decode date value=#{value}, format=#{I18n.t('date.formats.default')}"
        decoded_value = Date.strptime(value, I18n.t('date.formats.default'))
        log_debug "Row#write_value decoded value=#{decoded_value}"
      rescue
        log_debug "Row#write_value invalid date"
        errors.add(output, I18n.t('error.badformat', :column => column))
        decoded_value = nil
      end
    else
      decoded_value = value
    end
    send("#{output}=", decoded_value)
  end

  # Returns the workspace the row is attached to,
  # based on the existence of a reference column. 
  def workspace
    Workspace.select_entity_by_uuid(Workspace, self.workspace_uuid) if attribute_present?(:workspace_uuid)
  end

  def read_referenced_name_and_description(column, value)
    if value.present? 
      if column.kind == COLUMN_TYPE_REFERENCE and 
         column.grid_reference_uuid.present?
        log_debug "Row#read_referenced_name_and_description value=#{value}, grid_reference_uuid=#{column.grid_reference_uuid}"
        grid = column.grid_reference
        if grid.present? and grid.loaded
          log_debug "Row#read_referenced_name_and_description grid=#{grid.to_s}"
          return grid.select_reference_row_name_and_description(value)
        end
      end
      [value, ""]
    else
      ["", ""]
    end
  end

  # Indicates if attachments exist or not.
  def has_attachment?
    not self.attachments.nil? and not self.attachments.empty? 
  end

  # Returns the number of attachments.
  def count_attachments
    self.has_attachment? ? self.attachments.count : 0
  end

  # Removes the attachment that matches the given file name.
  def remove_attachment!(input_file)
    if input_file.present?
      log_debug "Row#remove_attachment!(#{input_file.original_filename})"
      for attachment in self.attachments
        if attachment.original_file_name.present? and
           attachment.original_file_name == input_file.original_filename
          log_debug "Row#remove_attachment! destroy"
          attachment.document = nil
          attachment.destroy
        end
      end
    end
  end

  # Creates a new associated locale row.
  def new_loc
    loc = row_locs.new
    loc.uuid = self.uuid
    loc.version = self.version
    loc
  end

  # Imports the instance of the row in the database.
  def import(xml_attribute, xml_value)
    log_debug "Row#import(#{xml_attribute}, #{xml_value})"
    case xml_attribute
      when 'grid_uuid' then self.grid_uuid = xml_value
      else super
    end
  end

  # Imports attribute value from the xml flow into the object.
  def import_attribute(xml_attribute, xml_value)
    log_debug "Row#import_attribute(#{xml_attribute}, #{xml_value})"
    self.initialization
    grid.load if grid.present? and not grid.loaded 
    grid.column_all.each do |column|
      log_debug "Row#import_attribute column=#{column}, column.uuid=#{column.uuid}"
      if xml_attribute == column.uuid
        write_value(column, xml_value)
      end
    end
  end

  # Copies attributes from the object to the target entity.
  def copy_attributes(entity)
    log_debug "Row#copy_attributes"
    super
    entity.grid_uuid = self.grid_uuid
    grid.load if grid.present? and not grid.loaded 
    grid.column_all.each do |column|
      log_debug "Row#copy_attributes column=#{column}"
      write_value(column, self.read_value(column))
    end
  end

  # Imports the instance of the object in the database,
  # as a new instance or as an update of an existing instance.
  def import!
    log_debug "Row#import!"
    if grid.nil?
      log_error "Row#import! Can't import row when there is no grid reference"
    else
      grid.load if grid.present? and not grid.loaded
      row = grid.row_select_entity_by_uuid_version(self.uuid, self.version)
      if row.present?
        log_debug "Row#import! present self=#{self.revision} row=#{row.revision}"
        if self.revision > row.revision 
          log_debug "Row#import! update"
          copy_attributes(row)
          row.update_user_uuid = Entity.session_user_uuid
          row.updated_at = Time.now
          make_audit(Audit::IMPORT)
          updated = grid.update_row!(self)
          grid.row_update_dates!(self.uuid)
          return "updated" if updated
        else
          log_debug "Row#import! skip update"
          return "skipped"
        end
      else
        log_debug "Row#import! new"
        self.create_user_uuid = self.update_user_uuid = Entity.session_user_uuid
        self.created_at = self.updated_at = Time.now
        make_audit(Audit::IMPORT)
        created = grid.create_row!(self)
        grid.row_update_dates!(self.uuid)
        return "inserted" if created
      end
    end
    ""
  end

  # Imports the given loc data into the appropriate locale row.
  # Fetches on the collection of local row and copies the attributes
  # of the provided loc data into the row that matches the language.
  def import_loc!(loc)
    log_debug "Row#import_loc!(loc=#{loc})"
    if grid.nil?
      log_error "Row#import_loc! Can't import row loc when there is no grid reference"
    else
      log_debug "Row#import_loc! grid=#{grid.inspect}"
      updated = 0
      grid.row_loc_select_entity_by_uuid(self.uuid, self.version).each do |row_loc|
        if row_loc.base_locale == loc.base_locale
          log_debug "Row#import_loc! update"
          loc.copy_attributes(row_loc)
          grid.update_row_loc!(self, row_loc)
          updated += 1
        end
      end
      if updated == 0
        log_debug "Row#import_loc! new"
        grid.create_row_loc!(self, loc)
      end
    end
  end

  # Creates local row for all the installed languages
  # that is not created yet for the given collection.
  # This insures on row exists for any installed language.
  def create_missing_loc!
    log_debug "Row#create_missing_loc!"
    if grid.nil?
      log_error "Row#create_missing_loc! Can't create row loc when there is no grid reference"
    else 
      base_locs = []
      base_loc = nil
      foundI18n = false
      grid.row_loc_select_entity_by_uuid(self.uuid, self.version).each do |loc|
        base_locs << loc.locale
        if not foundI18n
          base_loc = loc
          foundI18n = (loc.base_locale == I18n.locale.to_s)
        end
      end
      if base_loc.present?
        log_debug "Row#create_missing! base_locs=#{base_locs.inspect}"
        LANGUAGES.each do |lang, locale|
          if (base_locs.find {|value| locale.to_s == value}).nil?
            log_debug "Row#create_missing! locale=#{locale.to_s}"
            loc = new_loc
            base_loc.copy_attributes(loc)
            loc.locale = locale.to_s
            loc.base_locale = base_loc.base_locale
            log_debug "Row#create_missing! create_row_loc!"
            grid.create_row_loc!(self, loc)
          else
            log_debug "Row#create_missing! skip"
          end
        end
      end
    end
  end

private

  def self.loc_select_columns
    "row_locs.id, row_locs.uuid, row_locs.version, row_locs.lock_version, " + 
    "row_locs.base_locale, row_locs.locale, row_locs.name, row_locs.description"
  end
end