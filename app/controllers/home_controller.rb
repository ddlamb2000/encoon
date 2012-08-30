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
class HomeController < ApplicationController
  before_filter :load_workspaces, :only => [:index, 
                                            :refresh]

  def index
    @page_title = ""
    @page_icon = "home"
    unlock_as_of_date
    push_history
  end
  
  def refresh
    log_debug "HomeController#refresh date=#{params[:home][:session_date]}"
    session[:as_of_date] = Date.strptime(params[:home][:session_date], I18n.t('datepicker.decode'))
    redirect_to session[:last_url]
  end

  def history
    render :partial => "history"
  end

  # Sets the session flag referenced by the given parameter.
  def set
    session[params[:flag]] = true unless params[:flag].nil?
    render :nothing => true
  end

  # Unsets the session flag referenced by the given parameter.
  def unset
    session[params[:flag]] = false unless params[:flag].nil?
    render :nothing => true
  end
end