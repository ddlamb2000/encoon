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
  ROOT_UUID = 'a6beb050-0427-012d-09bf-0015c5091b01'
  ROOT_GRID_UUID = 'bdbcc440-0427-012d-09c5-0015c5091b01'
  ROOT_DB_TABLE = 'ce612c50-0427-012d-09c8-0015c5091b01'
  ROOT_DB_LOC_TABLE = 'd8badca0-0427-012d-09cb-0015c5091b01'
  XML_TAG = 'gridmapping'

  belongs_to :grid, :foreign_key => "grid_uuid", :primary_key => "uuid"
  validates_presence_of :grid_uuid, :db_table

  def import(xml_attribute, xml_value)
    log_debug "GridMapping#import(xml_attribute=#{xml_attribute}, " + 
              "xml_value=#{xml_value})"
    case xml_attribute
      when 'grid_uuid' then self.grid_uuid = xml_value
      when 'db_table' then self.db_table = xml_value
      when 'db_loc_table' then self.db_loc_table = xml_value
      else super
    end
  end
  
  def copy_attributes(entity)
    log_debug "GridMapping#copy_attributes"
    super
    entity.grid_uuid = self.grid_uuid
    entity.db_table = self.db_table
    entity.db_loc_table = self.db_loc_table
  end

  def import!
    log_debug "GridMapping#import!"
    log_error "Can't import grid mapping when " + 
          "there is no grid reference" if grid.nil?
    mapping = grid.mapping_select_entity_by_uuid_version(self.uuid, self.version)
    if mapping.present?
      if self.revision > mapping.revision 
        log_debug "GridMapping#import! update"
        copy_attributes(mapping)
        self.update_user_uuid = Grid.session_user_uuid    
        make_audit(Audit::IMPORT)
        mapping.save!
        mapping.update_dates!(grid.grid_mappings)
        return true
      else
        log_debug "GridMapping#import! skip update"
      end
    else
      log_debug "GridMapping#import! new"
      self.create_user_uuid = self.update_user_uuid = Grid.session_user_uuid    
      self.created_at = self.updated_at = Time.now    
      make_audit(Audit::IMPORT)
      save!
      update_dates!(grid.grid_mappings)
      return true
    end
    false
  end

  def create_missing_loc!
  end
end