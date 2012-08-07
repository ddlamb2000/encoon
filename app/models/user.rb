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
  ROOT_UUID = '08beadc0-ea31-012c-105d-00166f92f624'
  ROOT_EMAIL_UUID = 'e878a5b1-1a56-012d-50c9-4417fe7fde95'
  ROOT_FIRST_NAME_UUID = 'd164bfd1-1a56-012d-2c94-4417fe7fde95'
  ROOT_LAST_NAME_UUID = 'da06d511-1a56-012d-71ca-4417fe7fde95'
  
  SYSTEM_ADMINISTRATOR_UUID = 'eebdc1a0-dd45-012c-aafe-0026b0d63708'

  has_many :row_passwords, :foreign_key => "uuid", :primary_key => "uuid"
  has_many :row_attachments, :foreign_key => "uuid", :primary_key => "uuid"

  validates_presence_of :first_name, :last_name

  devise :database_authenticatable, 
         :registerable, 
         :recoverable, 
         :rememberable, 
         :trackable, 
         :validatable, 
         :confirmable, 
         :lockable, 
         :omniauthable

  attr_accessible :email, :first_name, :last_name, :password, :password_confirmation, :remember_me

  def to_s
    name = ""
    name = self.first_name if !self.first_name.blank?
    name = name + " " + self.last_name if !self.last_name.blank?
    name = self.email if name.blank?
    name
  end
  
  def has_photo?
    for attachment in row_attachments
      return true if attachment.photo?
    end
    false
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
    log_debug "User#import_attribute(xml_attribute=#{xml_attribute}, " + 
              "xml_value=#{xml_value})"
    case xml_attribute
      when ROOT_EMAIL_UUID then self.email = xml_value
      when ROOT_FIRST_NAME_UUID then self.first_name = xml_value
      when ROOT_LAST_NAME_UUID then self.last_name = xml_value
    end
  end
  
  def import!
    log_debug "User#import!"
    user = User.select_entity_by_uuid_version(User, self.uuid, self.version)
    if user.present?
      if self.revision > user.revision 
        log_debug "User#import! update"
        copy_attributes(user)
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
      self.password = rand(36**15).to_s(36);
      log_debug "User#import! generated password=#{self.password}"
      make_audit(Audit::IMPORT)
      save!
      update_dates!(User)
      return "inserted"
    end
    ""
  end

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
end