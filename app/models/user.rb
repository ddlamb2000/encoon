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
class User < Entity
  has_many :row_passwords, :foreign_key => "uuid", :primary_key => "uuid"
  has_many :attachments, :foreign_key => "uuid", :primary_key => "uuid"

  validates_presence_of :first_name, :last_name

  after_create :create_workspace_new_user

  devise :database_authenticatable, 
         :registerable, 
         :recoverable, 
         :rememberable, 
         :trackable, 
         :validatable, 
         :confirmable, 
         :lockable

  attr_accessible :email, :first_name, :last_name, :password, :password_confirmation, :remember_me

  def to_s
    name = ""
    name = self.first_name if !self.first_name.blank?
    name = name + " " + self.last_name if !self.last_name.blank?
    name = self.email if name.blank?
    name
  end
  
  def self.select_entity_by_uuid_version(collection, uuid, version)
    collection.find(:first, 
                    :select => self.all_select_columns,
                    :conditions => 
                      ["users.uuid = :uuid " + 
                       " AND users.version = :version ", 
                       {:uuid => uuid, 
                       :version => version}]) 
  end
  
  def import_attribute(xml_attribute, xml_value)
    log_debug "User#import_attribute(#{xml_attribute}, #{xml_value})"
    case xml_attribute
      when USER_EMAIL_UUID then self.email = xml_value
      when USER_FIRST_NAME_UUID then self.first_name = xml_value
      when USER_LAST_NAME_UUID then self.last_name = xml_value
    end
  end
  
  # Imports the instance of the object in the database,
  # as a new instance or as an update of an existing instance.
  def import!
    log_debug "User#import!"
    user = User.select_entity_by_uuid_version(User, self.uuid, self.version)
    if user.present?
      if self.updated_at > user.updated_at 
        log_debug "User#import! update"
        copy_attributes(user)
        user.update_user_uuid = Entity.session_user_uuid
        user.updated_at = Time.now
        make_audit(Audit::IMPORT)
        user.save!
        user.update_dates!(User)
        return "updated"
      else
        log_debug "User#import! skip update"
        return "skipped"
      end
    else
      log_debug "User#import! new"
      self.create_user_uuid = self.update_user_uuid = Entity.session_user_uuid
      self.created_at = self.updated_at = Time.now
      self.password = rand(36**15).to_s(36);
      log_debug "User#import! generated password=#{self.password}"
      make_audit(Audit::IMPORT)
      save!
      update_dates!(User)
      return "inserted"
    end
    ""
  end

  # Creates local row for all the installed languages: no locale for this class.
  def create_missing_loc! ; end

private

  def self.all_select_columns
    "users.id, users.uuid, " + 
    "users.version, users.lock_version, " + 
    "users.begin, users.end, users.enabled, " +
    "users.email, users.first_name, users.last_name, " +
    "users.created_at, users.updated_at, " +
    "users.create_user_uuid, users.update_user_uuid"
  end
  
  # Creates a private workspace for the user being created. 
  def create_workspace_new_user
    self.create_user_uuid = self.uuid
    save!
    log_debug "User#create_workspace_new_user workspace"
    workspace = Workspace.new
    workspace.create_user_uuid = self.uuid
    workspace.update_user_uuid = self.uuid
    workspace.public = false
    workspace.uri = self.email
    workspace.clean_uri!
    workspace.save!
    log_debug "User#create_workspace_new_user workspace_loc"
    workspace_loc = WorkspaceLoc.new
    workspace_loc.uuid = workspace.uuid
    workspace_loc.version = workspace.version
    workspace_loc.base_locale = I18n.locale.to_s
    workspace_loc.locale = I18n.locale.to_s
    workspace_loc.name = "#{self.first_name} #{self.last_name}"
    workspace_loc.save!
    log_debug "User#create_workspace_new_user missing loc"
    workspace.create_missing_loc!
    log_debug "User#create_workspace_new_user workspace sharing"
    workspace_sharing = WorkspaceSharing.new
    workspace_sharing.workspace_uuid = workspace.uuid
    workspace_sharing.user_uuid = self.uuid
    workspace_sharing.role_uuid = ROLE_TOTAL_CONTROL_UUID
    workspace_sharing.save!
    log_debug "User#create_workspace_new_user audit"
    workspace.make_audit(Audit::CREATE)
  end
end