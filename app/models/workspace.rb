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
#
# Workspaces are used to register sets of grids. 
# Every grid is attached to a workspace.
class Workspace < Entity
  has_many :grids, :foreign_key => "workspace_uuid", :primary_key => "uuid"
  has_many :workspace_locs,  :foreign_key => "uuid", :primary_key => "uuid"

  # Selects data based on its uuid in the given collection.
  def self.select_entity_by_uuid(collection, uuid)
    collection.
      select(self.all_select_columns).
      joins(:workspace_locs).
      where("workspaces.uuid = ?", uuid).
      where(as_of_date_clause("workspaces")).
      where(Grid::workspace_security_clause("workspaces")).
      where("workspace_locs.version = workspaces.version").
      where(locale_clause("workspace_locs")).
      first
  end

  # Selects data based on its uri in the given collection.
  def self.select_entity_by_uri(collection, uri)
    collection.
      select(self.all_select_columns).
      joins(:workspace_locs).
      where("workspaces.uri = ?", uri).
      where(as_of_date_clause("workspaces")).
      where(Grid::workspace_security_clause("workspaces")).
      where("workspace_locs.version = workspaces.version").
      where(locale_clause("workspace_locs")).
      first
  end

  # Selects data based on its uuid in the given collection and for a given version number.
  def self.select_entity_by_uuid_version(collection, uuid, version)
    collection.
      select(self.all_select_columns).
      joins(:workspace_locs).
      where("workspaces.uuid = ?", uuid).
      where("workspaces.version = ?", version).
      where(as_of_date_clause("workspaces")).
      where(Grid::workspace_security_clause("workspaces")).
      where("workspace_locs.version = workspaces.version").
      where(locale_clause("workspace_locs")).
      first
  end

  # Selects the workspaces the current user can access to.
  def self.user_workspaces
    select(self.user_select_columns).
    joins(:workspace_locs).
    where("workspace_locs.version = workspaces.version").
    where(as_of_date_clause("workspaces")).
    where(Grid::workspace_security_clause("workspaces", true, false)).
    where(locale_clause("workspace_locs")).
    order("workspace_locs.name")
  end

  # Creates a new associated locale row.
  def new_loc
    loc = workspace_locs.new
    loc.uuid = self.uuid
    loc.version = self.version
    loc
  end

  # Copies attributes from the object to the target entity.
  def copy_attributes(entity)
    log_debug "Workspace#copy_attributes"
    super
    entity.public = self.public
    entity.default_role_uuid = self.default_role_uuid
    entity.uri = self.uri
  end

  # Imports attribute value from the xml flow into the object.
  def import_attribute(xml_attribute, xml_value)
    log_debug "Workspace#import_attribute(#{xml_attribute}, #{xml_value})"
    case xml_attribute
      when WORKSPACE_PUBLIC_UUID then self.public = ['true','t','1'].include?(xml_value)
      when WORKSPACE_DEFAULT_ROLE_UUID then self.default_role_uuid = xml_value
    end
  end

  # Imports the instance of the object in the database,
  # as a new instance or as an update of an existing instance.
  def import!
    log_debug "Workspace#import!"
    workspace = Workspace.select_entity_by_uuid_version(Workspace, self.uuid, self.version)
    if workspace.present?
      if self.revision > workspace.revision 
        log_debug "Workspace#import! update"
        copy_attributes(workspace)
        workspace.update_user_uuid = Entity.session_user_uuid
        workspace.updated_at = Time.now
        make_audit(Audit::IMPORT)
        workspace.save!
        workspace.update_dates!(Workspace)
        return "updated"
      else
        log_debug "Workspace#import! skip update"
        return "skipped"
      end
    else
      log_debug "Workspace#import! new"
      self.create_user_uuid = self.update_user_uuid = Entity.session_user_uuid
      self.created_at = self.updated_at = Time.now
      make_audit(Audit::IMPORT)
      save!
      update_dates!(Workspace)
      return "inserted"
    end
    ""
  end

  # Imports the given loc data into the appropriate locale row.
  # Fetches on the collection of local row and copies the attributes
  # of the provided loc data into the row that matches the language.
  def import_loc!(loc)
    log_debug "Workspace#import_loc!"
    import_loc_base!(Workspace.locales(workspace_locs, self.uuid, self.version), loc)
  end

  # Creates local row for all the installed languages
  # that is not created yet for the given collection.
  # This insures on row exists for any installed language.
  def create_missing_loc!
    log_debug "Workspace#create_missing_loc!"
    create_missing_loc_base!(Workspace.locales(workspace_locs, self.uuid, self.version))
  end

private

  def self.user_select_columns
    "workspaces.id, workspaces.uuid, workspaces.uri, workspace_locs.name, workspace_locs.description"
  end

  def self.all_select_columns
    "workspaces.id, workspaces.uuid, workspaces.version, workspaces.lock_version, " + 
    "workspaces.begin, workspaces.end, workspaces.enabled, " +
    "workspaces.created_at, workspaces.updated_at, " +
    "workspaces.create_user_uuid, workspaces.update_user_uuid, " +
    "workspaces.public, workspaces.default_role_uuid, workspaces.uri, " +
    "workspace_locs.base_locale, workspace_locs.locale, workspace_locs.name, workspace_locs.description"
  end
end