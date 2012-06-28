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
# Entities are abstract active records identified using 
# an universal unique identifier,
# are version-based, dated and user-tracked. 
# Every entity is managed through a data grid.
class Entity < ActiveRecord::Base
  self.abstract_class = true
  
  @@uuid_gen = nil
  @@begin_of_time = Date::civil(1,1,1) 
  @@end_of_time = Date::civil(9999,12,31)

  belongs_to :create_user, 
             :select => "id, uuid, version, begin, end, enabled, email, first_name, last_name" , 
             :class_name => "User", 
             :foreign_key => "create_user_uuid", 
             :primary_key => "uuid"
  belongs_to :update_user, 
             :select => "id, uuid, version, begin, end, enabled, email, first_name, last_name" , 
             :class_name => "User", 
             :foreign_key => "update_user_uuid", 
             :primary_key => "uuid"
  has_many :audits, :foreign_key => "uuid", :primary_key => "uuid"

  validates_presence_of :uuid, :version
  validate :valid_dates
  
  before_validation(:on => :create) do
    defaults
  end

  def uuid_gen
    @@uuid_gen = UUID.new if @@uuid_gen.nil?
    @@uuid_gen
  end

  def self.begin_of_time
    @@begin_of_time
  end

  def self.end_of_time
    @@end_of_time
  end

  # Defaults begin date and version
  def defaults
    if self.uuid.blank?
      self.uuid = uuid_gen.generate
      self.enabled = true
    end
    self.version = 1 if self.version.blank?
    default_dates
  end
  
  def default_dates
    self.begin = @@begin_of_time if self.begin.blank? 
    self.end = @@end_of_time if self.end.blank?
  end
  
  def valid_dates
    if self.begin.present? and self.end.present?
      self.end = self.begin if self.end < self.begin
    end
  end
  
  def has_begin?
    self.begin.present? and self.begin != @@begin_of_time 
  end
 
  def has_end?
    self.end.present? and self.end != @@end_of_time 
  end
 
  # Returns the name of the user who created the record
  def who_created
    self.create_user.blank? ? "?" : self.create_user
  end

  # Returns the name of the user who updated the record
  def who_updated
    self.update_user.blank? ? "?" : self.update_user
  end
  
  def was_updated?
    self.updated_at != self.created_at
  end
  
  # Returns a revision number
  def revision
    1 + self.lock_version
  end

  # Returns true if the user who created the record has a photo
  def created_photo?
    if self.create_user.blank?
      false
    else
      self.create_user.has_photo?
    end
  end

  # Returns true if the user who updated the record has a photo
  def updated_photo?
    if self.update_user.blank?
      false
    else
      self.update_user.has_photo?
    end
  end

  def to_s
    name
  end

  def name
    attribute_present?(:name) ? read_attribute(:name) : ""
  end

  def description
    attribute_present?(:description) ? read_attribute(:description) : ""
  end

  # Indicates if the given uuid is a valid uuid
  def self.uuid?(uuid)
    uuid.present? and uuid.length == 36
  end

  def self.loc_select_columns
    "id, uuid, version, lock_version, " +
    "base_locale, locale, " +
    "name, description"
  end
  
  def self.all_locales(collection, uuid, version)
    log_debug "Entity#all_locales(uuid=#{uuid}), version=#{version.to_s}"
    collection.find(:all, 
                    :select => self.loc_select_columns,
                    :conditions => 
                           ["version = :version AND locale = base_locale", 
                           {:version => version}], 
                    :order => "locale")
  end

  def self.as_of_date_clause(synonym)
    "'#{session_as_of_date}' BETWEEN #{synonym}.begin AND #{synonym}.end" + 
    " AND #{synonym}.enabled = 't'"
  end
  
  def as_of_date_clause(synonym)
    Entity.as_of_date_clause(synonym) 
  end
  
  def self.locale_clause(synonym)
    "#{synonym}.locale = '#{session_locale}'"
  end

  def locale_clause(synonym)
    Entity.locale_clause(synonym)
  end

  def export(xml)
    log_debug "Entity#export"
    xml.uuid(self.uuid)
    xml.version(self.version)
    xml.begin(self.begin) if has_begin?
    xml.end(self.end) if has_end?
    xml.revision(self.revision)
    xml.enabled(self.enabled)
    xml.created_at(self.created_at)
    xml.updated_at(self.updated_at) if was_updated?
    xml.create_user_uuid(self.create_user_uuid)
    xml.update_user_uuid(self.update_user_uuid) if was_updated?
  end

  def import(xml_attribute, xml_value)
    log_debug "Entity#import(xml_attribute=#{xml_attribute}, " + 
              "xml_value=#{xml_value})"
    case xml_attribute
      when 'uuid' then self.uuid = xml_value
      when 'version' then self.version = xml_value.to_i
      when 'begin' then self.begin = Date::parse(xml_value)
      when 'end' then self.end = Date::parse(xml_value)
      when 'enabled' then self.enabled = (xml_value == 'true')
      when 'revision' then self.lock_version = xml_value.to_i-1
    end
  end
  
  def import_attribute(xml_attribute, xml_value)
    log_debug "Entity#import_attribute(xml_attribute=#{xml_attribute}, " + 
              "xml_value=#{xml_value})"
  end
  
  def copy_attributes(entity)
    log_debug "Entity#copy_attributes"
    entity.begin = self.begin    
    entity.end = self.end
    entity.enabled = self.enabled
    entity.updated_at = self.updated_at    
    entity.update_user_uuid = self.update_user_uuid    
  end

  def import_loc_base!(collection, loc)
    log_debug "Entity#import_loc_base!(loc=#{loc})"
    updated = 0
    collection.each do |entity_loc|
      if entity_loc.base_locale == loc.base_locale
        log_debug "Entity#import_loc_base! update"
        loc.copy_attributes(entity_loc)
        entity_loc.save!
        updated += 1
      end
    end
    if updated == 0
      log_debug "Entity#import_loc_base! new"
      loc.save!
    end
  end

  def create_missing_loc_base!(collection)
    log_debug "Entity#create_missing_loc!"
    base_locs = []
    base_loc = nil
    foundI18n = false
    collection.each do |loc|
      base_locs << loc.base_locale
      if not foundI18n
        base_loc = loc
        foundI18n = (loc.base_locale == I18n.locale.to_s)
      end
    end
    if base_loc.present?
      LANGUAGES.each do |lang, locale|
        if (base_locs.find {|value| locale.to_s == value}).nil?
          loc = new_loc        
          base_loc.copy_attributes(loc)
          loc.locale = locale.to_s
          loc.base_locale = base_loc.base_locale
          loc.save!
          log_debug "Entity#create_missing_loc! new, locale=#{locale.to_s}"
        end
      end
    end
  end

  def all_audits
    log_debug "Entity#all_audits"
    audits.find(:all, :order => "updated_at")
  end

  def make_audit(kind)
    log_debug "Entity#make_audit(kind=#{kind})"
    audit = Audit.new
    audit.grid_uuid = self.grid_uuid if attribute_present?(:grid_uuid)
    audit.uuid = self.uuid
    audit.version = self.version
    audit.kind = kind
    audit.locale = I18n.locale.to_s
    audit.lock_version = self.lock_version
    audit.update_user_uuid = self.update_user_uuid
    audit.save!
  end

  def update_dates!(collection)
    log_debug "Entity#update_dates"
    previous_item = nil
    last_item = nil
    collection.find(:all, 
                    :conditions => 
                          ["uuid = :uuid", 
                          {:uuid => self.uuid}],
                    :order => "begin").each do |item|
      if item.enabled
        log_debug "Entity#update_dates previous_item=#{previous_item.inspect}"
        if previous_item.present? and previous_item.end != item.begin-1
          log_debug "Entity#update_dates previous_item set end date"
          previous_item.end = item.begin-1
          previous_item.update_user_uuid = Entity.session_user_uuid
          previous_item.save!
        end
        previous_item = item
      end
      last_item = item
    end
    log_debug "Entity#update_dates last_item=#{last_item.inspect}"
    if last_item.present? and last_item.end != @@end_of_time
      log_debug "Entity#update_dates last_item set end date"
      last_item.end = @@end_of_time
      last_item.update_user_uuid = Entity.session_user_uuid
      last_item.save!
    end
  end
  
  def self.session_user_uuid=user_uuid
    Thread.current[:session_user_uuid] = user_uuid
  end
  
  def self.session_user_uuid
    Thread.current[:session_user_uuid]
  end
  
  def self.session_user_display_name=user_display_name
    Thread.current[:session_user_display_name] = user_display_name
  end
  
  def self.session_user_display_name
    Thread.current[:session_user_display_name]
  end
  
  def self.session_as_of_date=as_of_date
    Thread.current[:session_as_of_date] = as_of_date
  end

  def self.session_as_of_date
    Thread.current[:session_as_of_date]
  end

  def self.session_locale=locale
    Thread.current[:session_locale] = locale
  end

  def self.session_locale
    Thread.current[:session_locale]
  end

  def self.workspace_security_clause(synonym)
    if DATA_GRID_SECURITY_ACTIVATED
      "(" +
      " EXISTS (" +
      "  SELECT 1 FROM workspace_sharings " +
      "  WHERE workspace_sharings.workspace_uuid = #{synonym}.uuid" + 
      "  AND " + as_of_date_clause("workspace_sharings") +
      "  AND workspace_sharings.user_uuid = '#{Entity.session_user_uuid}'" +
      "  AND workspace_sharings.role_uuid in ('#{Role::ROLE_READ_ONLY_UUID}', '#{Role::ROLE_READ_WRITE_UUID}', '#{Role::ROLE_READ_WRITE_ALL_UUID}', '#{Role::ROLE_TOTAL_CONTROL_UUID}')" +
      " )" +
      " OR EXISTS (" +
      "  SELECT 1 FROM workspaces workspace_security " +
      "  WHERE workspace_security.uuid = #{synonym}.uuid" + 
      "  AND " + as_of_date_clause("workspace_security") +
      "  AND (" +
      "   workspace_security.public = 't'" +
      "   OR workspace_security.create_user_uuid = '#{Entity.session_user_uuid}'" +
      "  )" +
      " )" +
      ")"
    else
      "1=1"
    end
  end

  def self.grid_security_clause(synonym)
    if DATA_GRID_SECURITY_ACTIVATED
      "(" +
      " EXISTS (" +
      "  SELECT 1 FROM grids grid_security, workspace_sharings " +
      "  WHERE grid_security.uuid = #{synonym}.uuid" + 
      "  AND workspace_sharings.workspace_uuid = grid_security.workspace_uuid" + 
      "  AND " + as_of_date_clause("grid_security") +
      "  AND " + as_of_date_clause("workspace_sharings") +
      "  AND workspace_sharings.user_uuid = '#{Entity.session_user_uuid}'" +
      "  AND workspace_sharings.role_uuid in ('#{Role::ROLE_READ_ONLY_UUID}', '#{Role::ROLE_READ_WRITE_UUID}', '#{Role::ROLE_READ_WRITE_ALL_UUID}', '#{Role::ROLE_TOTAL_CONTROL_UUID}')" +
      " )" + 
      " OR EXISTS (" +
      "  SELECT 1 FROM grids grid_security, workspaces workspace_security" +
      "  WHERE grid_security.uuid = #{synonym}.uuid" + 
      "  AND workspace_security.uuid = grid_security.workspace_uuid" + 
      "  AND " + as_of_date_clause("grid_security") +
      "  AND " + as_of_date_clause("workspace_security") +
      "  AND (" +
      "   workspace_security.public = 't'" +
      "   OR workspace_security.default_role_uuid in ('#{Role::ROLE_READ_ONLY_UUID}', '#{Role::ROLE_READ_WRITE_UUID}', '#{Role::ROLE_READ_WRITE_ALL_UUID}', '#{Role::ROLE_TOTAL_CONTROL_UUID}')" +
      "   OR workspace_security.create_user_uuid = '#{Entity.session_user_uuid}'" +
      "  )" +
      " )" +
      ")"
    else
      "1=1"
    end
  end

  # Logs debugging message.
  # Information related to the connected user,
  # as of date and language is logged automatically.
  #
  # ==== Parameters
  #
  # * +message+ - A string that is logged as debugging information.
  #
  # ==== Examples
  #
  #   # Prints a value in the log for debugging.
  #   Entity.log_debug "Entity#calc Calculated value = #{value.to_s}"
  def self.log_debug(message)
    logger.debug "[" +
                 "#{session_user_display_name}:" +
                 "#{I18n.l(session_as_of_date)}" +
                 "(#{session_locale})] " +
                 message
  end

  def self.log_warning(message)
    logger.warn "[" +
                 "#{session_user_display_name}:" +
                 "#{I18n.l(session_as_of_date)}" +
                 "(#{session_locale})] " +
                 "## WARNING ## " +
                 message
  end

  def self.log_error(message)
    logger.error "[" +
                 "#{session_user_display_name}:" +
                 "#{I18n.l(session_as_of_date)}" +
                 "(#{session_locale})] " +
                 "## ERROR ## " +
                 message
  end

  def log_debug(message) ; Entity.log_debug(message) ; end
  def log_warning(message) ; Entity.log_warning(message) ; end
  def log_error(message) ; Entity.log_error(message) ; end
end