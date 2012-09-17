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
class WorkspaceSharing < Entity
  ROOT_UUID = '4c2b6200-9a7a-012f-642b-4417fe7fde95'
  ROOT_WORKSPACE_UUID_UUID = '7e67ae50-9a7a-012f-642b-4417fe7fde95'
  ROOT_USER_UUID_UUID = '8c13aa60-9a7a-012f-642b-4417fe7fde95'
  ROOT_ROLE_UUID_UUID = '95d604f0-9c87-012f-089e-4417fe7fde95'

  has_many :grids, :foreign_key => "workspace_uuid", :primary_key => "uuid"
  
  def self.select_entity_by_uuid_version(collection, uuid, version)
    collection.find(:first, 
                    :select => self.all_select_columns,
                    :conditions => 
                      ["workspace_sharings.uuid = :uuid " + 
                       " AND workspace_sharings.version = :version ", 
                       {:uuid => uuid, 
                       :version => version}]) 
  end
  
  def copy_attributes(entity)
    log_debug "WorkspaceSharing#copy_attributes"
    super
    entity.workspace_uuid = self.workspace_uuid    
    entity.role_uuid = self.role_uuid    
    entity.user_uuid = self.user_uuid    
  end

  def import_attribute(xml_attribute, xml_value)
    log_debug "WorkspaceSharing#import_attribute(xml_attribute=#{xml_attribute}, " + 
              "xml_value=#{xml_value})"
    case xml_attribute
      when ROOT_WORKSPACE_UUID_UUID then self.workspace_uuid = xml_value
      when ROOT_USER_UUID_UUID then self.user_uuid = xml_value
      when ROOT_ROLE_UUID_UUID then self.role_uuid = xml_value
    end
  end

  def import!
    log_debug "WorkspaceSharing#import!"
    workspace = WorkspaceSharing.select_entity_by_uuid_version(WorkspaceSharing, 
                                                        self.uuid, 
                                                        self.version)
    if workspace.present?
      if self.revision > workspace.revision 
        log_debug "WorkspaceSharing#import! update"
        copy_attributes(workspace)
        self.update_user_uuid = Entity.session_user_uuid
        self.updated_at = Time.now
        make_audit(Audit::IMPORT)
        workspace.save!
        workspace.update_dates!(Workspace)
        return "updated"
      else
        log_debug "WorkspaceSharing#import! skip update"
        return "skipped"
      end
    else
      log_debug "WorkspaceSharing#import! new"
      self.create_user_uuid = self.update_user_uuid = Entity.session_user_uuid
      self.created_at = self.updated_at = Time.now
      make_audit(Audit::IMPORT)
      save!
      update_dates!(Workspace)
      return "inserted"
    end
    ""
  end

  def create_missing_loc! ; end

private

  def self.all_select_columns
    "workspace_sharings.id, workspace_sharings.uuid, " + 
    "workspace_sharings.version, workspace_sharings.lock_version, " + 
    "workspace_sharings.begin, workspace_sharings.end, workspace_sharings.enabled, " +
    "workspace_sharings.created_at, workspace_sharings.updated_at, " +
    "workspace_sharings.create_user_uuid, workspace_sharings.update_user_uuid, " +
    "workspace_sharings.workspace_uuid, workspace_sharings.user_uuid, workspace_sharings.role_uuid"
  end
end