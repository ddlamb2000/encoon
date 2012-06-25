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
  XML_TAG = 'workspace'

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
                       " AND " + security_clause("workspaces"), 
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
                       " AND " + security_clause("workspaces"), 
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
  
  def export(xml)
    log_debug "Workspace#export"
    xml.workspace(:title => self.name) do
      super(xml)
      Workspace.all_locales(workspace_locs, 
                            self.uuid, 
                            self.version).each do |loc|
        log_debug "Workspace#export locale #{loc.base_locale}"
        xml.locale do
          xml.base_locale(loc.base_locale)
          xml.locale(loc.locale)
          xml.name(loc.name)
          xml.description(loc.description) if loc.description.present?
        end
      end
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
        self.update_user_uuid = Grid.session_user_uuid    
        make_audit(Audit::IMPORT)
        workspace.save!
        workspace.update_dates!(Workspace)
        return true
      else
        log_debug "Workspace#import! skip update"
      end
    else
      log_debug "Workspace#import! new"
      self.create_user_uuid = self.update_user_uuid = Grid.session_user_uuid    
      self.created_at = self.updated_at = Time.now    
      make_audit(Audit::IMPORT)
      save!
      update_dates!(Workspace)
      return true
    end
    false
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
    "workspace_locs.base_locale, workspace_locs.locale, " +
    "workspace_locs.name, workspace_locs.description"
  end
  
  def self.security_clause(synonym)
    if DATA_GRID_SECURITY_ACTIVATED
      "(" +
      #" #{synonym}.create_user_uuid = '#{session_user_uuid}'" +
      #" OR " +
      " EXISTS (" +
      "  SELECT 1 FROM workspace_sharings " +
      "  WHERE workspace_sharings.workspace_uuid = #{synonym}.uuid" + 
      "  AND " + as_of_date_clause("workspace_sharings") +
      "  AND workspace_sharings.user_uuid = '#{session_user_uuid}'" +
      "  AND workspace_sharings.role_uuid in ('#{Role::ROLE_READ_ONLY_UUID}', '#{Role::ROLE_READ_WRITE_UUID}', '#{Role::ROLE_TOTAL_CONTROL_UUID}')" +
      " )" +
      ")"
    else
      "1=1"
    end
  end
end

class WorkspaceLoc < EntityLoc
  validates_presence_of :name
end

class WorkspaceSharing < Entity
  ROOT_UUID = '4c2b6200-9a7a-012f-642b-4417fe7fde95'
end