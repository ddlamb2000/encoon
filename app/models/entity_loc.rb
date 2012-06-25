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
# Use case
# installed languages: en fr es
# 
# user connected with locale=fr
# creation of a new item, name "Voitur"
# ==> locale  base_locale  name
#     en      fr           Voitur    (new)
#     fr      fr           Voitur    (new)
#     es      fr           Voitur    (new)
# ==> locale  base_locale  name
#     fr      fr           Voitur    (selected)
#
# update of item, name "Voiture" -- update when locale = fr or base_locale = fr
# ==> locale  base_locale  name
#     en      fr           Voiture   (updated)
#     fr      fr           Voiture   (updated)
#     es      fr           Voiture   (updated)
# 
# user connected with locale=en
# ==> locale  base_locale  name
#     en      fr           Voiture   (selected)
#
# update of item, name "Car" -- update when locale = en or base_locale = en
# ==> locale  base_locale  name
#     en      en           Car       (updated)
#     fr      fr           Voiture
#     es      fr           Voiture
#
class EntityLoc < ActiveRecord::Base
  self.abstract_class = true
  validates_presence_of :locale

  def to_s
    attribute_present?(:name) ? read_attribute(:name) : ""
  end

  def name
    attribute_present?(:name) ? read_attribute(:name) : ""
  end

  def description
    attribute_present?(:description) ? read_attribute(:description) : ""
  end

  def import(xml_attribute, xml_value)
    log_debug "EntityLoc::import(xml_attribute=#{xml_attribute}, " + 
              "xml_value=#{xml_value})"
    case xml_attribute
      when 'base_locale' then self.base_locale = xml_value
      when 'locale' then self.locale = xml_value
      when 'name' then self.name = xml_value
      when 'description' then self.description = xml_value
    end
  end

  def copy_attributes(entity_loc)
    log_debug "EntityLoc::copy_attributes entity_loc=#{entity_loc.inspect}"
    entity_loc.base_locale = self.base_locale    
    entity_loc.locale = self.locale
    entity_loc.name = self.name
    entity_loc.description = self.description    
  end
end