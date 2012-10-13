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
class ColumnMapping < Entity
  belongs_to :column, :foreign_key => "column_uuid", :primary_key => "uuid"
  validates_presence_of :column_uuid, :db_column

  def import_attribute(xml_attribute, xml_value)
    log_debug "ColumnMapping#import_attribute(#{xml_attribute}, #{xml_value})"
    case xml_attribute
      when COLUMN_MAPPING_COLUMN_UUID then self.column_uuid = xml_value
      when COLUMN_MAPPING_DB_COLUMN then self.db_column = xml_value
    end
  end
  
  def copy_attributes(entity)
    log_debug "ColumnMapping#copy_attributes"
    super
    entity.column_uuid = self.column_uuid
    entity.db_column = self.db_column
  end

  def import!
    log_debug "ColumnMapping#import!"
    log_error "Can't import column mapping when " + 
          "there is no column reference" if column.nil?
    mapping = column.column_mapping_select_entity_by_uuid_version(self.uuid, 
                                                                  self.version)
    if mapping.present?
      if self.revision > mapping.revision 
        log_debug "ColumnMapping#import! update"
        copy_attributes(mapping)
        self.update_user_uuid = Entity.session_user_uuid
        self.updated_at = Time.now
        make_audit(Audit::IMPORT)
        mapping.save!
        mapping.update_dates!(column.column_mappings)
        return "updated"
      else
        log_debug "ColumnMapping#import! skip update"
        return "skipped"
      end
    else
      log_debug "ColumnMapping#import! new"
      self.create_user_uuid = self.update_user_uuid = Entity.session_user_uuid
      self.created_at = self.updated_at = Time.now
      make_audit(Audit::IMPORT)
      save!
      update_dates!(column.column_mappings)
      return "inserted"
    end
    ""
  end

  # Creates local row for all the installed languages: no locale for this class.
  def create_missing_loc! ; end
end