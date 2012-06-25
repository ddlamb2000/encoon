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
class Audit < ActiveRecord::Base
  CREATE = "CREATE"
  UPDATE = "UPDATE"
  DELETE = "DELETE"
  IMPORT = "IMPORT"
  ATTACH = "ATTACH"
  DETACH = "DETACH"
  PASSWORD = "PASSWORD"

  belongs_to :update_user, 
             :select => "id, uuid, " + 
                        "email, first_name, last_name, " + 
                        "version, begin, end, enabled", 
             :class_name => "User", 
             :foreign_key => "update_user_uuid", 
             :primary_key => "uuid"
  validates_presence_of :kind
  
  def kind_name
    language = LANGUAGES.find {|lang, locale| self.locale == locale}
    case kind
      when CREATE then 
        I18n.t('general.audit_create',
            :version => self.version,
            :time => self.updated_at,
            :language => language[0])
      when UPDATE then 
        I18n.t('general.audit_update', 
            :version => self.version,
            :time => self.updated_at, 
            :revision => self.lock_version+1,
            :language => language[0])
      when DELETE then 
        I18n.t('general.audit_delete',
            :version => self.version,
            :time => self.updated_at,
            :language => language[0])
      when ATTACH then 
        I18n.t('general.audit_attach')
      when DETACH then 
        I18n.t('general.audit_detach')
      when PASSWORD then 
        I18n.t('general.audit_pass')
      else ""
    end
  end
      
  # Returns the name of the user who updated the record
  def who_updated
    self.update_user
  end
end