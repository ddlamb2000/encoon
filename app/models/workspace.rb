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
class Workspace < Entity
  ROOT_UUID = 'f54309d0-ea30-012c-1054-00166f92f624'
  ROOT_DEFAULT_ROLE_UUID_UUID = '0e1f4990-a26c-012f-de85-4417fe7fde95'
  ROOT_PUBLIC_UUID = 'f1ca9820-a26b-012f-de85-4417fe7fde95'
  ROOT_URI_UUID = '44b03b60-a26c-012f-de85-4417fe7fde95'

  has_many :grids, :foreign_key => "workspace_uuid", :primary_key => "uuid"
  has_many :workspace_locs,  :foreign_key => "uuid", :primary_key => "uuid"
  
  def before_destroy
    log_debug "Workspace#before_destroy"
    super
    WorkspaceLoc.destroy_all(["uuid = :uuid AND version = :version", 
                             {:uuid => self.uuid, :version => self.version}])
    if Workspace.all_versions(Workspace, self.uuid).length == 0
      log_debug "Workspace#before_destroy remove_orphans"
      grids.destroy_all
    end
  end
  
  # Selects data based on uuid
  def self.select_entity_by_uuid(collection, uuid)
    collection.find(:first, 
                    :joins => :workspace_locs,
                    :select => self.all_select_columns,
                    :conditions => 
                      ["workspaces.uuid = :uuid " + 
                       " AND " + as_of_date_clause("workspaces") +
                       " AND workspace_locs.version = workspaces.version " +
                       " AND " + locale_clause("workspace_locs"), 
                       {:uuid => uuid}]) 
  end
  
  def self.select_entity_by_uuid_version(collection, uuid, version)
    collection.find(:first, 
                    :joins => :workspace_locs,
                    :select => self.all_select_columns,
                    :conditions => 
                      ["workspaces.uuid = :uuid " + 
                       " AND workspaces.version = :version " + 
                       " AND workspace_locs.version = workspaces.version " +
                       " AND " + locale_clause("workspace_locs"), 
                       {:uuid => uuid, 
                       :version => version}]) 
  end
  
  def self.all_versions(collection, uuid)
    collection.find(:all, 
                    :joins => :workspace_locs,
                    :select => self.all_select_columns,
                    :conditions => 
                      ["workspaces.uuid = :uuid " +
                       " AND workspace_locs.version = workspaces.version " + 
                       " AND " + security_clause("workspaces"), 
                       " AND " + locale_clause("workspace_locs"), 
                       {:uuid => uuid}], 
                    :order => "workspaces.begin")
  end
  
  def new_loc
    loc = workspace_locs.new
    loc.uuid = self.uuid
    loc.version = self.version
    loc
  end
  
  def copy_attributes(entity)
    log_debug "Workspace#copy_attributes"
    super
    entity.public = self.public    
    entity.default_role_uuid = self.default_role_uuid    
    entity.uri = self.uri    
  end

  def import_attribute(xml_attribute, xml_value)
    log_debug "Workspace#import_attribute(xml_attribute=#{xml_attribute}, " + 
              "xml_value=#{xml_value})"
    case xml_attribute
      when ROOT_PUBLIC_UUID then self.public = ['true','t','1'].include?(xml_value)
      when ROOT_DEFAULT_ROLE_UUID_UUID then self.default_role_uuid = xml_value
      when ROOT_URI_UUID then self.uri = xml_value
    end
  end

  def import!
    log_debug "Workspace#import!"
    workspace = Workspace.select_entity_by_uuid_version(Workspace, 
                                                        self.uuid, 
                                                        self.version)
    if workspace.present?
      if self.revision > workspace.revision 
        log_debug "Workspace#import! update"
        copy_attributes(workspace)
        make_audit(Audit::IMPORT)
        workspace.save!
        workspace.update_dates!(Workspace)
        return "updated"
      else
        log_debug "Workspace#import! skip update"
      end
    else
      log_debug "Workspace#import! new"
      make_audit(Audit::IMPORT)
      save!
      update_dates!(Workspace)
      return "inserted"
    end
    ""
  end

  def import_loc!(loc)
    import_loc_base!(Workspace.all_locales(workspace_locs, 
                                           self.uuid, 
                                           self.version), loc)
  end

  def create_missing_loc!
    create_missing_loc_base!(Workspace.all_locales(workspace_locs, 
                                                   self.uuid, 
                                                   self.version))
  end
  
private

  def self.all_select_columns
    "workspaces.id, workspaces.uuid, " + 
    "workspaces.version, workspaces.lock_version, " + 
    "workspaces.begin, workspaces.end, workspaces.enabled, " +
    "workspaces.created_at, workspaces.updated_at, " +
    "workspaces.create_user_uuid, workspaces.update_user_uuid, " +
    "workspaces.public, workspaces.default_role_uuid, workspaces.uri, " +
    "workspace_locs.base_locale, workspace_locs.locale, " +
    "workspace_locs.name, workspace_locs.description"
  end
end

class WorkspaceLoc < EntityLoc
  validates_presence_of :name
end