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
  def upload(file_name)
    self.file_name = file_name
    self.data_file=File.new(self.file_name)
    save!
  end
  
  def data_file=(input_file)
    require "rexml/document"
    if input_file.present? 
      self.file_name = input_file.original_filename if self.file_name.blank?
      xml_doc = REXML::Document.new input_file.read
      xml_doc.elements.each do |xml_root| 
        xml_root.elements.each do |xml_entity| 
          log_debug "Upload#data_file= xml_entity=#{xml_entity.name}"
          entity = nil
          locs = []
          case xml_entity.name
            when Workspace::XML_TAG then entity = Workspace.new
            when Grid::XML_TAG then entity = Grid.new
            when Column::XML_TAG then entity = Column.new
            when Row::XML_TAG then entity = Row.new
            when User::XML_TAG then entity = User.new
            when GridMapping::XML_TAG then entity = GridMapping.new
            when ColumnMapping::XML_TAG then entity = ColumnMapping.new
            else log_warning "Unkown xml tag #{xml_entity.name}"
          end
          entity.transaction do
            xml_entity.elements.each do |xml_attribute|
              if xml_attribute.name == 'locale' and 
                 xml_attribute.has_elements?
                entity_loc = entity.new_loc
                xml_attribute.elements.each do |xml_attribute_loc|
                  entity_loc.import(xml_attribute_loc.name, undecode(xml_attribute_loc))      
                end
                locs << entity_loc
              elsif xml_attribute.name == 'data' and 
                    xml_attribute.has_attributes?
                entity.import_attribute(xml_attribute.attributes['uuid'], undecode(xml_attribute))      
              else
                entity.import(xml_attribute.name, 
                              undecode(xml_attribute))   
              end
            end
            entity.default_dates
            if entity.import!
              if locs.length > 0
                locs.each do |entity_loc|
                  entity.import_loc!(entity_loc)
                end
              end
              entity.create_missing_loc!
            end
          end
        end
      end
    end
  end

private
  
  def undecode(attribute)
    attribute.get_text.to_s.gsub(/&amp;/,'&').gsub(/&lt;/,'<').gsub(/&gt;/,'>')    
  end
end