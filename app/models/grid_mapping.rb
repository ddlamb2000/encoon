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
class GridMapping < Entity
  belongs_to :grid, :foreign_key => "grid_uuid", :primary_key => "uuid"
  validates_presence_of :grid_uuid, :db_table

  def import_attribute(xml_attribute, xml_value)
    log_debug "GridMapping#import_attribute(#{xml_attribute}, #{xml_value})"
    case xml_attribute
      when GRID_MAPPING_GRID_UUID then self.grid_uuid = xml_value
      when GRID_MAPPING_DB_TABLE then self.db_table = xml_value
      when GRID_MAPPING_DB_LOC_TABLE then self.db_loc_table = xml_value
    end
  end
  
  def copy_attributes(entity)
    log_debug "GridMapping#copy_attributes"
    super
    entity.grid_uuid = self.grid_uuid
    entity.db_table = self.db_table
    entity.db_loc_table = self.db_loc_table
  end

  # Imports the instance of the object in the database,
  # as a new instance or as an update of an existing instance.
  def import!
    log_debug "GridMapping#import!"
    if grid.nil?
      log_error "Can't import grid mapping when there is no grid reference"
    else
      mapping = grid.mapping_select_entity_by_uuid_version(self.uuid, self.version)
      if mapping.present?
        if self.revision > mapping.revision 
          log_debug "GridMapping#import! update"
          copy_attributes(mapping)
          mapping.update_user_uuid = Entity.session_user_uuid
          mapping.updated_at = Time.now
          make_audit(Audit::IMPORT)
          mapping.save!
          mapping.update_dates!(grid.grid_mappings)
          return "updated"
        else
          log_debug "GridMapping#import! skip update"
          return "skipped"
        end
      else
        log_debug "GridMapping#import! new"
        self.create_user_uuid = self.update_user_uuid = Entity.session_user_uuid
        self.created_at = self.updated_at = Time.now
        make_audit(Audit::IMPORT)
        save!
        update_dates!(grid.grid_mappings)
        return "inserted"
      end
    end
    ""
  end

  # Creates local row for all the installed languages: no locale for this class.
  def create_missing_loc! ; end
end