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
class Upload < Entity
  require "rexml/document"

  def upload(file_name)
    self.file_name = file_name
    self.data_file=File.new(self.file_name)
    save!
  end
  
  def data_file=(input_file)
    if not input_file.blank?
      self.records = self.inserted = self.updated = self.skipped = self.elapsed = 0
      start_run = Time.now 
      self.file_name = input_file.original_filename if self.file_name.blank?
      log_debug "Upload#data_file= [#{self.file_name}]"
      xml_doc = REXML::Document.new input_file.read
      xml_doc.elements.each do |xml_root| 
        log_debug "Upload#data_file= xml_root=#{xml_root.name}"
        if "encoon" == xml_root.name
          xml_root.elements.each do |xml_entity|
            log_debug "Upload#data_file= xml_entity=#{xml_entity.name}"
            if "row" == xml_entity.name
              self.records = self.records + 1
              grid_uuid = ""
              xml_entity.attributes.each do |xml_attribute|
                log_debug "Upload#data_file= xml_attribute=#{xml_attribute}"
                grid_uuid = xml_attribute[1] if "grid_uuid" == xml_attribute[0]
              end
              if grid_uuid == ""
                log_warning "Upload#data_file= Missing grid_uui"
              else
                entity = nil
                locs = []
                case grid_uuid
                  when WORKSPACE_UUID then entity = Workspace.new
                  when GRID_UUID then entity = Grid.new
                  when COLUMN_UUID then entity = Column.new
                  when USER_UUID then entity = User.new
                  when GRID_MAPPING_UUID then entity = GridMapping.new
                  when COLUMN_MAPPING_UUID then entity = ColumnMapping.new
                  when WORKSPACE_SHARING_UUID then entity = WorkspaceSharing.new
                  else entity = Row.new
                end
                entity.transaction do
                  xml_entity.elements.each do |xml_attribute|
                    if xml_attribute.name == 'locale' and xml_attribute.has_elements?
                      entity_loc = entity.new_loc
                      xml_attribute.elements.each do |xml_attribute_loc|
                        entity_loc.import(xml_attribute_loc.name, undecode(xml_attribute_loc))      
                      end
                      locs << entity_loc
                    elsif xml_attribute.name == 'data' and xml_attribute.has_attributes?
                      entity.import_attribute(xml_attribute.attributes['uuid'], undecode(xml_attribute))      
                    else
                      entity.import(xml_attribute.name, undecode(xml_attribute))   
                    end
                  end
                  entity.default_dates
                  imported = entity.import!
                  if imported != ""
                    if imported != "skipped"
                      locs.each { |entity_loc|  entity.import_loc!(entity_loc) } if locs.length > 0
                      entity.create_missing_loc!
                    end
                    self.inserted = self.inserted + 1 if "inserted" == imported
                    self.updated = self.updated + 1 if "updated" == imported
                    self.skipped = self.skipped + 1 if "skipped" == imported
                  end
                end
              end
            end
          end
        end
      end
      self.elapsed = 1000*(Time.now - start_run)
    end
  end

private

  def undecode(attribute)
    attribute.get_text.to_s.gsub(/&amp;/,'&').gsub(/&lt;/,'<').gsub(/&gt;/,'>').gsub(/&quot;/,'"')    
  end
end