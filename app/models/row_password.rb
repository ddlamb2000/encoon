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
class RowPassword < ActiveRecord::Base
  validate :password_confirmation
  validate :password_non_blank

  # User_password (transient attribute)
  def user_password
    @user_password
  end

  # Password is encrypted based on user entry
  def user_password=(pwd)
    @user_password = pwd
    return if pwd.blank?
    create_new_salt
    self.password = encrypted_password(self.user_password, self.salt)
  end
  
  def user_password_confirmation
    @user_password_confirmation
  end

  def user_password_confirmation=(pwd)
    @user_password_confirmation = pwd
  end
  
private
  
  # Controls the password is not blank
  def password_non_blank
    errors.add(:user_password, 
               I18n.t('error.mis_password')) if self.password.blank? 
  end
  
  # Controls the password confirmation
  def password_confirmation
    errors.add(:user_password_confirmation, 
               I18n.t('error.pass_conf')) if self.password.present? and 
                                             self.user_password_confirmation != 
                                             self.user_password   
  end
  
  # Encrypts a password
  def encrypted_password(user_password, salt)
    require 'digest/sha1'
    string_to_hash = user_password + "factory" + salt
    Digest::SHA1.hexdigest(string_to_hash)
  end
  
  # Create salt for password encryption
  def create_new_salt
    self.salt = self.object_id.to_s + rand.to_s
  end
end