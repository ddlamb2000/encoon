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
  XML_TAG = 'user'

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

  def export(xml)
    log_debug "User#export"
    xml.user(:title => self.to_s) do
      super(xml)
      xml.email(self.email)
      xml.firstname(self.first_name)
      xml.lastname(self.last_name)
    end
  end
  
  def import(xml_attribute, xml_value)
    log_debug "User#import(xml_attribute=#{xml_attribute}, " + 
              "xml_value=#{xml_value})"
    case xml_attribute
      when 'email' then self.email = xml_value
      when 'firstname' then self.first_name = xml_value
      when 'lastname' then self.last_name = xml_value
      else super
    end
  end
  
  def import!
    log_debug "User#import!"
    user = User.select_entity_by_uuid_version(User, self.uuid, self.version)
    if user.present?
      if self.revision > user.revision 
        log_debug "User#import! update"
        copy_attributes(user)
        self.update_user_uuid = Grid.session_user_uuid    
        make_audit(Audit::IMPORT)
        user.save!
        user.update_dates!(User)
        return true
      else
        log_debug "User#import! skip update"
      end
    else
      log_debug "User#import! new"
      self.create_user_uuid = self.update_user_uuid = Grid.session_user_uuid    
      self.created_at = self.updated_at = Time.now    
      make_audit(Audit::IMPORT)
      save!
      update_dates!(User)
      return true
    end
    false
  end

  def create_missing_loc!
  end
end