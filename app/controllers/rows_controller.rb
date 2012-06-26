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
class RowsController < ApplicationController
  before_filter :authenticate_user!
  before_filter :findParent
  before_filter :findEntity, :only => [:show, 
                                       :show_more, 
                                       :edit_inline, 
                                       :edit_row_inline, 
                                       :update, 
                                       :destroy,
                                       :attach_document,
                                       :save_attachment,
                                       :delete_attachment,
                                       :enter_password,
                                       :save_password,
                                       :photo,
                                       :file]

  def index
    log_debug "RowController#index: params=#{params.inspect}"
    if params[:format].nil? or params[:format] != 'xml'  
      set_page_title
      push_history
      unlock_as_of_date
    end
  end
  
  def search
    log_debug "RowController#search: params=#{params.inspect}"
    render :layout => false
  end
  
  def set_page_title
    if @grid.present?
      if @row.present?
        @page_title = I18n.t('general.object_name', :type => @grid, :name => @grid.row_title(@row))
      else
        @page_title = @grid.to_s
      end
      if @grid.uuid == Grid::ROOT_UUID
        @page_icon = "table"
      elsif @grid.uuid == Workspace::ROOT_UUID
        @page_icon = "workspace"
      else
        @page_icon = "entity"
      end
    else
      @page_title = I18n.t('error.no_grid')
      @page_icon = "exclamation"
    end
  end

  def show
    log_debug "RowController#show: params=#{params.inspect}"
    if params[:format].nil? or params[:format] != 'xml'  
      set_page_title
      push_history
    end
  end

  def show_more
    log_debug "RowController#show_more: params=#{params.inspect}"
    respond_to do |format|
      format.js
    end
  end

  def hide_more
    respond_to do |format|
      format.js
    end
  end

  def new_inline
    log_debug "RowController#new_inline: params=#{params.inspect}"
    @filters = params[:filters]
    @filters_uuid = get_filters_uuid(@filters)
    @grid.load_cached_grid_structure(@filters, true)
    @row = @grid.rows.build
    @row_loc = RowLoc.new
    @grid.row_initialization(@row, @filters)
    set_page_title
    respond_to do |format|
      format.js
    end
  end

  def edit_inline
    log_debug "RowController#edit_inline: params=#{params.inspect}"
    set_page_title
    respond_to do |format|
      format.js
    end
  end

  def edit_row_inline
    log_debug "RowController#edit_row_inline: params=#{params.inspect}"
    set_page_title
    respond_to do |format|
      format.js
    end
  end

  def create
    log_debug "RowController#create: params=#{params.inspect}"
    saved = false
    @grid.load_cached_grid_structure(params[:filters], true)
    @row = @grid.rows.new
    @row.initialization
    begin
      @row.transaction do
        @row.begin = param_begin_date
        @row.create_user_uuid = current_user.uuid
        @row.update_user_uuid = current_user.uuid
        populate_from_params
        @grid.load_cached_grid_structure(params[:filters], false)
        if @grid.has_translation?
          LANGUAGES.each do |lang, locale|
            @row_loc = @row.new_loc
            @row_loc.uuid = @row.uuid
            @row_loc.version = @row.version
            @row_loc.base_locale = I18n.locale.to_s
            @row_loc.locale = locale
            @row_loc.name = params[:name]
            @row_loc.description = params[:description]
            if @grid.row_loc_validate(@row, @row_loc, Grid::PHASE_CREATE)
              @grid.create_row_loc!(@row_loc)
            else
              @grid.row_validate(@row, Grid::PHASE_CREATE)
              log_debug "RowController#create: rollback!"
              raise ActiveRecord::Rollback
            end
          end
        end
        if @grid.row_validate(@row, Grid::PHASE_CREATE)
          @grid.create_row!(@row)
          saved = true
        else
          log_debug "RowController#create: rollback!(2)"
          raise ActiveRecord::Rollback
        end
      end
    rescue ActiveRecord::RecordInvalid => invalid
      log_debug "RowController#create: invalid=#{invalid.inspect}"
      saved = false
    end
    change_as_of_date(@row)
    respond_to do |format|
      if saved
        @row = @grid.row_select_entity_by_uuid(@row.uuid)
        name = @grid.row_title(@row)
        flash[:notice] = I18n.t('transaction.created', 
                                :type => @grid, 
                                :name => name)
        format.html { redirect_to session[:last_url] }
      else
        log_debug "RowController#create: error, params=#{params.inspect}"
        set_page_title
        format.html { render :action => "_new_inline" }
      end
    end
  end
  
  def update
    log_debug "RowController#update: params=#{params.inspect}"
    saved = false
    @grid.load_cached_grid_structure(params[:filters])
    begin
      @grid.transaction do
        version = @row.version
        if params[:mode] == Grid::PHASE_NEW_VERSION then
          new_version = @grid.row_max_version(@row.uuid)+1
          @row = @row.clone
          @row.version = new_version
          @row.lock_version = 0
          @row.begin = param_begin_date
          if multiple_versions_ok?
            @row.enabled = params[:row][:enabled]
            @row.create_user_uuid = current_user.uuid
            @row.update_user_uuid = current_user.uuid
            populate_from_params
            if @grid.row_validate(@row, Grid::PHASE_NEW_VERSION)
              @grid.create_row!(@row)
              if @grid.has_translation?
                for @row_loc in @grid.row_loc_select_entity_by_uuid(@row.uuid, version)
                  @row_loc = @row_loc.clone
                  @row_loc.version = new_version
                  @row_loc.lock_version = 0
                  if @row_loc.locale == I18n.locale.to_s or 
                     @row_loc.base_locale == I18n.locale.to_s
                    @row_loc.base_locale = I18n.locale.to_s
                    @row_loc.name = params[:name]
                    @row_loc.description = params[:description]
                  end
                  if @grid.row_loc_validate(@row, 
                                            @row_loc, 
                                            Grid::PHASE_NEW_VERSION)
                    @grid.create_row_loc!(@row_loc)
                  else
                    @grid.row_validate(@row, Grid::PHASE_UPDATE)
                    log_debug "RowController#update: rollback!"
                    raise ActiveRecord::Rollback
                  end
                end
              end
              saved = true
            else
              raise ActiveRecord::Rollback
            end
          end
        else
          @row.begin = param_begin_date
          @row.update_user_uuid = current_user.uuid
          if multiple_versions_ok?
            @row.enabled = params[:row][:enabled]
            populate_from_params
            if @grid.row_validate(@row, Grid::PHASE_UPDATE)
              @grid.update_row!(@row)
              if @grid.has_translation?
                for @row_loc in @grid.row_loc_select_entity_by_uuid(@row.uuid, 
                                                                    version)
                  if @row_loc.locale == I18n.locale.to_s or 
                     @row_loc.base_locale == I18n.locale.to_s
                    @row_loc.base_locale = I18n.locale.to_s
                    @row_loc.name = params[:name]
                    @row_loc.description = params[:description]
                    if @grid.row_loc_validate(@row, 
                                              @row_loc, 
                                              Grid::PHASE_UPDATE)
                      @grid.update_row_loc!(@row_loc)
                    else
                      log_debug "RowController#update: rollback!(3)"
                      raise ActiveRecord::Rollback
                    end
                  end
                end
              end
              saved = true
            else
              log_debug "RowController#update: rollback!(4)"
              raise ActiveRecord::Rollback
            end
          end
        end
      end
      @grid.row_update_dates!(@row.uuid)
    rescue ActiveRecord::RecordInvalid => invalid
      log_debug "RowController#update: invalid=#{invalid.inspect}"
      saved = false
    end
    change_as_of_date(@row)
    respond_to do |format|
      if saved
        name = @grid.row_title(@row)
        flash[:notice] = I18n.t('transaction.updated', 
                                :type => @grid, :name => name)
        format.html { redirect_to session[:last_url] }
      else
        log_debug "RowController#update: error, params=#{params.inspect}"
        set_page_title
        format.html { render :action => "_edit_inline" }
      end
    end
  end

  def destroy
    log_debug "RowController#destroy: params=#{params.inspect}"
    saved = false
    name = @grid.row_title(@row)
    begin
      @row.transaction do
        @row.enabled = false
        if @grid.row_validate(@row, Grid::PHASE_DESTROY)
          @grid.update_row!(@row)
          @grid.destroy_row!(@row)
          @grid.row_update_dates!(@row.uuid)
          if @grid.has_translation?
            for @row_loc in @grid.row_loc_select_entity_by_uuid(@row.uuid)
              @grid.destroy_row_loc!(@row_loc)
            end
          end
          saved = true
        end
      end
    rescue ActiveRecord::RecordInvalid => invalid
      log_debug "RowController#destroy: invalid=#{invalid.inspect}"
      saved = false
    end
    respond_to do |format|
      if saved
        flash[:notice] = I18n.t('transaction.deleted', 
                                :type => @grid, :name => name)
        format.html { redirect_to session[params[:inline] ? :last_url : :prior_url] }
      else
        log_debug "RowController#destroy: error, params=#{params.inspect}"
        format.html { render :action => "show" }
      end
    end
  end

  def attach_document
    log_debug "RowController#attach_document: params=#{params.inspect}"
    @row_attachment = @row.row_attachments.new
    set_page_title 
    respond_to do |format|
      format.js
    end
  end

  def save_attachment
    log_debug "RowController#save_attachment: params=#{params.inspect}"
    saved = false
    begin
      @row.transaction do
        @row.remove_attachment!(params[:row_attachment][:attach_document])
        @row_attachment = @row.row_attachments.new
        @row_attachment.update_attributes(params[:row_attachment])
        @row.make_audit(Audit::ATTACH)
        saved = true
      end
    rescue ActiveRecord::RecordInvalid => invalid
      log_debug "RowController#save_attachment: invalid=#{invalid.inspect}"
      saved = false
    end
    respond_to do |format|
      if saved
        name = @grid.row_title(@row)
        flash[:notice] = I18n.t('transaction.attached', 
                                :type => @grid, :name => name)
        format.html { redirect_to session[:last_url] }
      else
        log_debug "RowController#save_document: error"
        format.html { render :action => "attach_document" }
      end
    end
  end

  def enter_password
    logg_debug "RowController#enter_password: params=#{params.inspect}"
    @row_password = @row.row_passwords.new 
    render :layout => false
  end

  def save_password
    log_debug "RowController#save_password: params=#{params.inspect}"
    saved = false
    begin
      @row.transaction do
        @row.remove_password!
        @row_password = @row.row_passwords.new
        @row_password.uuid = @row.uuid
        @row_password.user_password = 
          params[:row_password][:user_password]
        @row_password.user_password_confirmation = 
          params[:row_password][:user_password_confirmation]
        @row_password.save!
        @row.make_audit(Audit::PASSWORD)
        saved = true
      end
    rescue ActiveRecord::RecordInvalid => invalid
      log_debug "RowController#save_password: invalid=#{invalid.inspect}"
      saved = false
    end
    respond_to do |format|
      if saved 
        name = @grid.row_title(@row)
        flash[:notice] = I18n.t('transaction.pass_created', :name => name)
        format.html { redirect_to session[:last_url] }
      else
        log_debug "RowController#save_password: error"
        format.html { render :action => "enter_password" }
      end
    end
  end

  def photo
    log_debug "RowController#photo: params=#{params.inspect}"
    if params[:photo_id].present?
      @row_attachment = @row.row_attachments.find(params[:photo_id])
    else
      @row_attachment = @row.first_photo
    end
    if @row_attachment.present?
      send_data(@row_attachment.document)
    end
  end

  def file
    log_debug "RowController#file: params=#{params.inspect}"
    @row_attachment = @row.row_attachments.find(params[:file_id])
    if @row_attachment.present?
      send_data @row_attachment.document, 
                :type => @row_attachment.content_type,
                :filename => @row_attachment.file_name
    end
  end

  def delete_attachment
    log_debug "RowController#delete_attachment: params=#{params.inspect}"
    saved = false
    begin
      @row.transaction do
        @row_attachment = @row.row_attachments.find(params[:file_id])
        if @row_attachment.present?
          @row_attachment.delete
          @row.make_audit(Audit::DETACH)
        end
      end
      saved = true
    rescue ActiveRecord::RecordInvalid => invalid
      log_debug "RowController#delete_attachment: invalid=#{invalid.inspect}"
      saved = false
    end
    if saved
      flash[:notice] = I18n.t('transaction.deleted', 
                              :type => @row_attachment.content_type, 
                              :name => @row_attachment.file_name)
    end
    respond_to do |format|
      format.html { render :action => "show" }
    end
  end

private

  def populate_from_params
    for column in @grid.column_all
      @row.write_value(column, params[:row][column.physical_column])
    end
  end

  # This is used to generate a unique id in the page
  # in order to avoid duplicates.
  # This is based on the concatenation of identifiers used to filter data 
  def get_filters_uuid(filters)
    output = ""
    if filters.present?
      filters.collect do |filter|
         output << filter[:column_uuid]
      end
    end
    output
  end

  def multiple_versions_ok?
    if params[:row][:enabled] == "0" and 
        !@grid.row_enabled_version_exists?(@row, 
                                           params[:mode] == 'new_version')
      @row.errors.add_to_base(I18n.t('error.ena_record'))
      return false
    elsif @grid.row_begin_duplicate_exists?(@row, 
                                            params[:begin_date])
      @row.errors.add_to_base(I18n.t('error.dup_record', 
                                     :date => params[:begin_date]))
      return false
    end
    true
  end
  
  def findParent
    @grid = nil
    if Entity.uuid?(params[:grid_id])
      @grid = Grid.select_entity_by_uuid(Grid, params[:grid_id])
    else
      @grid = Grid.select_entity_by_id(Grid, params[:grid_id])
    end
    if @grid.nil?
      Entity.log_debug "RowsController#findParent " + 
                       "Invalid: can't find data grid #{params[:grid_id]}"
    else
      @grid.load_cached_grid_structure    
    end
  end

  def findEntity
    if @grid.present?
      if Entity.uuid?(params[:id])
        @row = @row_loc = @grid.row_select_entity_by_uuid(params[:id])
        unlock_as_of_date
      else
        @row = @row_loc = @grid.row_select_entity_by_id(params[:id])
        lock_as_of_date
      end
      Entity.log_debug "RowsController#findEntity " + 
                       "Invalid: can't find row with " +
                       "grid_id=#{params[:grid_id].to_s}" +
                       " and id=#{params[:id].to_s}" if @row.nil?
      change_as_of_date(@row) if @row.present?
    end
  end
end