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
class GridController < ApplicationController
  before_filter :load_workspaces, :only => [:home, :show, :refresh]

  before_filter :authenticate_user!, :only => [:new,
                                               :edit,
                                               :create,
                                               :update,
                                               :attributes,
                                               :details,
                                               :attach,
                                               :save_attachment,
                                               :delete_attachment,
                                               :import,
                                               :upload]

  # Renders the home page using hard-coded references.
  def home
    params[:workspace] = Workspace::SYSTEM_WORKSPACE_URI
    params[:grid] = Grid::HOME_PAGE_UUID
    params[:row] = Grid::HOME_WELCOME_UUID
    show
  end

  # Renders credits using hard-coded references.
  def credits
    params[:workspace] = Workspace::SYSTEM_WORKSPACE_URI
    params[:grid] = Grid::HOME_PAGE_UUID
    params[:row] = Grid::HOME_CREDITS_UUID
    show
  end

  def history
    session[:show_history] = true
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

  def refresh
    log_debug "GridController#refresh date=#{params[:home][:session_date]}"
    session[:as_of_date] = Date.strptime(params[:home][:session_date], t('datepicker.decode'))
    redirect_to session[:last_url]
  end

  def show
    log_debug "GridController#show"
    @filters = params[:filters]
    selectWorkspaceAndGrid
    if @grid.present?
      @grid.load_cached_grid_structure(@filters)
      selectRow
      @table_columns = @grid.filtered_columns
      if @grid.uuid == Grid::ROOT_UUID
        if @row.present?
          @grid_cast = @grid.select_grid_cast(@row.uuid)
          @attached_grids = [@grid_cast]
        end
      else
        @attached_grids = Grid.select_referenced_grids(@grid.uuid) if @grid.present?
      end
    end
    respond_to do |format|
      format.html do
        set_page_title
        push_history
        render :show, :status => @status
      end
      format.xml do
        if @grid.uuid == Grid::ROOT_UUID and @grid_cast.present?
            @table_rows = @grid_cast.row_all(@filters, nil, -1, true)
        end
        render :show, :status => @status
      end
    end
  end

  # Renders the content of a list through an Ajax request  
  def list
    log_debug "GridController#list"
    @filters = params[:filters]
    @search = params[:search]
    @page = params[:page]
    selectGridAndWorkspace
    if @grid.present?
      @grid.load_cached_grid_structure(@filters) 
      @table_row_count = @grid.row_count(@filters)
      @table_rows = @grid.row_all(@filters, @search, @page, true)
      @table_columns = @grid.filtered_columns
    end
    render :partial => "list"
  end
  
  # Renders the details of an article through an Ajax request  
  def details
    log_debug "GridController#details"
    selectGridAndWorkspace
    @grid.load_cached_grid_structure if @grid.present?
    selectRow
    if @grid.present? and @row.present?
      @versions = @grid.row_all_versions(@row.uuid)
      @locales = @grid.row_all_locales(@row.uuid, @row.version)
      @audits = @row.all_audits
    end
    session[:show_details] = true
    render :partial => "details"
  end
  
  # Renders the attachments of an article through an Ajax request  
  def attachments
    log_debug "GridController#attachments"
    selectGridAndWorkspace
    @grid.load_cached_grid_structure if @grid.present?
    selectRow
    render :partial => "attachments"
  end
  
  def attributes
    log_debug "GridController#attributes"
    selectGridAndWorkspace
    @grid.load_cached_grid_structure if @grid.present?
    @columns = @grid.column_all
    render :partial => "attributes"
  end

  def import
    log_debug "GridController#import"
    selectGridAndWorkspace
    @grid.load_cached_grid_structure if @grid.present?
    render :partial => "import"
  end
  
  def upload
    log_debug "GridController#upload"
    saved = false
    selectGridAndWorkspace
    @grid.load_cached_grid_structure if @grid.present?
    @upload = Upload.new
    @upload.create_user_uuid = @upload.update_user_uuid = Entity.session_user_uuid
    begin
      @upload.transaction do
        log_debug "GridController#upload: initialize transaction " +
                  "params[:data_file]=#{params[:data_file]}"
        if params[:data_file].blank?
          @upload.errors.add(:file_name, I18n.t('error.required', :column => I18n.t('field.file')))
        else
          @upload.update_attributes(:data_file => params[:data_file])
          @upload.save!
          saved = true
        end
      end
    rescue ActiveRecord::RecordInvalid => invalid
      log_debug "GridController#upload: invalid=#{invalid.inspect}"
      saved = false
    rescue Exception => invalid
      log_error "GridController#upload", invalid
      saved = false
    end
    respond_to do |format|
      if saved
        log_debug "GridController#upload: saved"
        flash[:notice] = I18n.t('message.uploaded',
                                :record_count => @upload.records,
                                :insert_count => @upload.inserted,
                                :update_count => @upload.updated)
        render :nothing => true
      else
        log_debug "GridController#upload: error, @upload.errors=#{@upload.errors.inspect}"
        format.html { render :json => @upload.errors, :status => :unprocessable_entity }
      end
    end
  end

  # Renders the details of an article through an Ajax request  
  def row
    log_debug "GridController#row"
    @filters = params[:filters]
    selectGridAndWorkspace
    @grid.load_cached_grid_structure
    selectRow
    @table_columns = @grid.filtered_columns
    set_page_title
    render :partial => "row"
  end

  def new
    log_debug "GridController#new"
    @container = params[:container]
    @refresh_list = params[:refresh_list]
    @filters = params[:filters]
    selectGridAndWorkspace
    @grid.load_cached_grid_structure(@filters) if @grid.present?
    @filters_uuid = get_filters_uuid(@filters)
    @row = @grid.rows.build
    @row_loc = RowLoc.new
    @grid.row_initialization(@row, @filters)
    render :partial => "edit", :locals => {:new_row => true}
  end

  def edit
    log_debug "GridController#edit"
    @container = params[:container]
    @refresh_list = params[:refresh_list]
    @filters = params[:filters]
    selectGridAndWorkspace
    @grid.load_cached_grid_structure(@filters) if @grid.present?
    selectRow
    render :partial => "edit", :locals => {:new_row => false}
  end

  def create
    log_debug "GridController#create"
    saved = false
    @filters = params[:filters]
    selectGridAndWorkspace
    @grid.load_cached_grid_structure(@filters, true) if @grid.present?
    @row = @grid.rows.new
    @row.initialization
    begin
      @row.transaction do
        log_debug "GridController#create: initialize transaction"
        @row.begin = param_begin_date
        log_debug "GridController#create: populate from parameter values"
        populate_from_params
        @grid.load_cached_grid_structure(@filters)
        if @grid.has_translation?
          LANGUAGES.each do |lang, locale|
            log_debug "GridController#create: locale=#{locale}"
            @row_loc = @row.new_loc
            @row_loc.uuid = @row.uuid
            @row_loc.version = @row.version
            @row_loc.base_locale = I18n.locale.to_s
            @row_loc.locale = locale
            @row_loc.name = params[:name]
            @row_loc.description = params[:description]
            log_debug "GridController#create: row_loc_validate"
            if @grid.row_loc_validate(@row, @row_loc, Grid::PHASE_CREATE)
              log_debug "GridController#create: create_row_loc!"
              @grid.create_row_loc!(@row_loc)
            else
              @grid.row_validate(@row, Grid::PHASE_CREATE)
              log_debug "GridController#create: rollback!"
              raise ActiveRecord::Rollback
            end
          end
        end
        log_debug "GridController#create: row_validate"
        if @grid.row_validate(@row, Grid::PHASE_CREATE)
          log_debug "GridController#create: create_row!"
          @grid.create_row!(@row)
          saved = true
        else
          log_debug "GridController#create: rollback!(2)"
          raise ActiveRecord::Rollback
        end
      end
    rescue ActiveRecord::RecordInvalid => invalid
      log_debug "GridController#create: invalid=#{invalid.inspect}"
      saved = false
    rescue Exception => invalid
      log_error "GridController#create", invalid
      saved = false
    end
    change_as_of_date(@row)
    respond_to do |format|
      if saved
        format.html {list}
      else
        log_debug "GridController#create: error, @row.errors=#{@row.errors.inspect}"
        format.html { render :json => @row.errors, :status => :unprocessable_entity }
      end
    end
  end
  
  def update
    log_debug "GridController#update"
    saved = false
    @refresh_list = params[:refresh_list]
    @filters = params[:filters]
    selectGridAndWorkspace
    @grid.load_cached_grid_structure(@filters) if @grid.present?
    selectRow
    if params[:lock_version] != @row.lock_version.to_s
      log_debug "GridController#update locked! " +
                "params[:lock_version]=#{params[:lock_version]}, " +
                "@row.lock_version=#{@row.lock_version}"
      @row.errors.add(:uuid, I18n.t('error.locked'))
    else
      begin
        @grid.transaction do
          log_debug "GridController#update initialize transaction"
          version = @row.version
          if params[:mode] == Grid::PHASE_NEW_VERSION then
            log_debug "GridController#update new version"
            new_version = @grid.row_max_version(@row.uuid)+1
            @row = @row.clone
            @row.version = new_version
            @row.lock_version = 0
            @row.begin = param_begin_date
            if multiple_versions_ok?
              @row.enabled = params[:enabled]
              log_debug "GridController#update populate from parameter values"
              populate_from_params
              log_debug "GridController#update row_validate"
              if @grid.row_validate(@row, Grid::PHASE_NEW_VERSION)
                log_debug "GridController#update create_row!"
                @grid.create_row!(@row)
                if @grid.has_translation?
                  for @row_loc in @grid.row_loc_select_entity_by_uuid(@row.uuid, version)
                    log_debug "GridController#update update row_loc values"
                    @row_loc = @row_loc.clone
                    @row_loc.version = new_version
                    @row_loc.lock_version = 0
                    log_debug "GridController#update locale=#{@row_loc.locale}"
                    if @row_loc.locale == I18n.locale.to_s or 
                       @row_loc.base_locale == I18n.locale.to_s
                      log_debug "GridController#update update row_loc values"
                      @row_loc.base_locale = I18n.locale.to_s
                      @row_loc.name = params[:name]
                      @row_loc.description = params[:description]
                    end
                    log_debug "GridController#update row_loc_validate"
                    if @grid.row_loc_validate(@row, 
                                              @row_loc, 
                                              Grid::PHASE_NEW_VERSION)
                      log_debug "GridController#update create_row_loc!"
                      @grid.create_row_loc!(@row_loc)
                    else
                      log_debug "GridController#update row_validate"
                      @grid.row_validate(@row, Grid::PHASE_UPDATE)
                      log_debug "GridController#update: rollback!"
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
            log_debug "GridController#update update existing version"
            @row.begin = param_begin_date
            if multiple_versions_ok?
              log_debug "GridController#update versions OK"
              @row.enabled = params[:enabled]
              log_debug "GridController#update populate from parameter values"
              populate_from_params
              log_debug "GridController#update row_validate"
              if @grid.row_validate(@row, Grid::PHASE_UPDATE)
                log_debug "GridController#update update_row!"
                @grid.update_row!(@row)
                if @grid.has_translation?
                  for @row_loc in @grid.row_loc_select_entity_by_uuid(@row.uuid, version)
                    log_debug "GridController#update locale=#{@row_loc.locale}"
                    if @row_loc.locale == I18n.locale.to_s or 
                       @row_loc.base_locale == I18n.locale.to_s
                      log_debug "GridController#update update row_loc values"
                      @row_loc.base_locale = I18n.locale.to_s
                      @row_loc.name = params[:name]
                      @row_loc.description = params[:description]
                      log_debug "GridController#update row_loc_validate"
                      if @grid.row_loc_validate(@row, 
                                                @row_loc, 
                                                Grid::PHASE_UPDATE)
                        log_debug "GridController#update update_row_loc!"
                        @grid.update_row_loc!(@row_loc)
                      else
                        log_debug "GridController#update: rollback!(3)"
                        raise ActiveRecord::Rollback
                      end
                    end
                  end
                end
                saved = true
              else
                log_debug "GridController#update: rollback!(4)"
                raise ActiveRecord::Rollback
              end
            end
          end
        end
        log_debug "GridController#update row_update_dates!"
        @grid.row_update_dates!(@row.uuid)
      rescue ActiveRecord::RecordInvalid => invalid
        log_debug "GridController#update: invalid=#{invalid.inspect}"
        saved = false
      rescue Exception => invalid
        log_error "GridController#update", invalid
        saved = false
      end
    end
    change_as_of_date(@row)
    respond_to do |format|
      if saved
        if @container == ""
          render :nothing => true
        else
          format.html {@refresh_list ? list : row}
        end
      else
        log_debug "GridController#update: error, params=#{params.inspect}"
        format.html { render :json => @row.errors, :status => :unprocessable_entity }
      end
    end
  end

  def attach
    log_debug "GridController#attach"
    selectGridAndWorkspace
    @grid.load_cached_grid_structure if @grid.present?
    selectRow
    @attachment = @row.attachments.new
    render :partial => "attach"
  end
  
  def save_attachment
    log_debug "GridController#save_attachment"
    selectGridAndWorkspace
    @grid.load_cached_grid_structure if @grid.present?
    selectRow
    saved = false
    @attachment = @row.attachments.new
    begin
      if params[:document].blank?
        @attachment.errors.add(:document_file_name, I18n.t('error.required', :column => I18n.t('field.file')))
      else
        @row.transaction do
          @row.update_user_uuid = Entity.session_user_uuid
          @row.remove_attachment!(params[:document])
          @attachment.original_file_name = (params[:document]).original_filename
          @attachment.document = params[:document]
          @attachment.create_user_uuid = Entity.session_user_uuid
          @attachment.save!
          @row.save!
          @row.make_audit(Audit::ATTACH)
          saved = true
        end
      end
    rescue ActiveRecord::RecordInvalid => invalid
      log_debug "GridController#save_attachment: invalid=#{invalid.inspect}"
      saved = false
    rescue Exception => invalid
      log_error "GridController#save_attachment", invalid
      saved = false
    end
    respond_to do |format|
      if saved
        render :text => "<div id='ok'>OK</div>"
        return
      else
        log_debug "GridController#save_document: error"
        format.html { render :json => @attachment.errors, :status => :unprocessable_entity }
      end
    end
  end

  def delete_attachment
    log_debug "GridController#delete_attachment"
    selectGridAndWorkspace
    @grid.load_cached_grid_structure
    selectRow
    saved = false
    begin
      @row.transaction do
        @attachment = @row.attachments.find(params[:id])
        if @attachment.present?
          @attachment.delete
          @row.make_audit(Audit::DETACH)
        end
      end
      saved = true
    rescue ActiveRecord::RecordInvalid => invalid
      log_debug "GridController#delete_attachment: invalid=#{invalid.inspect}"
      saved = false
    rescue Exception => invalid
      log_error "GridController#delete_attachment", invalid
      saved = false
    end
    row
  end

private

  def set_page_title
    @status = 200
    unless @grid.nil?
      unless @row.nil?
        @page_title = t('general.object_name', :type => @grid, :name => @grid.row_title(@row))
      else
        @status = 404
        @page_title = t('error.no_data')
      end
    else
      @status = 404
      @page_title = t('error.no_grid')
    end
  end

  def populate_from_params
    for column in @grid.column_all
      value = params["row_#{column.physical_column}"]
      log_debug "GridController#populate_from_params" +
                " column.physical_column=#{column.physical_column}" +
                " value=#{value}"
      @row.write_value(column, value)
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
    if params[:enabled] == "0" and 
        !@grid.row_enabled_version_exists?(@row, params[:mode] == 'new_version')
      @row.errors.add_to_base(t('error.ena_record'))
      return false
    elsif @grid.row_begin_duplicate_exists?(@row, params[:begin_date])
      @row.errors.add_to_base(t('error.dup_record', :date => params[:begin_date]))
      return false
    end
    true
  end
  
  def selectWorkspaceAndGrid
    Entity.log_debug "GridController#selectWorkspaceAndGrid"
    @workspace = nil
    @grid = nil
    if params[:workspace].present?
      if params[:grid].nil? or params[:row].nil?
        parameters = params[:workspace].split("/")
        Entity.log_debug "GridController#selectWorkspaceAndGrid parameters=#{parameters.inspect}"
        params[:workspace] = parameters[0]
        params[:grid] = parameters[1]
        params[:row] = parameters[2]
        if params[:grid].nil? and params[:row].nil?
          Entity.log_debug "GridController#selectWorkspaceAndGrid workspace only"
          params[:row] = params[:workspace]
          params[:workspace] = Workspace::SYSTEM_WORKSPACE_URI
          params[:grid] = Workspace::ROOT_UUID
        elsif not params[:grid].nil? and params[:row].nil?
          Entity.log_debug "GridController#selectWorkspaceAndGrid grid only"
          params[:row] = params[:grid]
          params[:workspace] = Workspace::SYSTEM_WORKSPACE_URI
          params[:grid] = Grid::ROOT_UUID
        end
      end 
      if Entity.uuid?(params[:workspace])
        @workspace = Workspace.select_entity_by_uuid(Workspace, params[:workspace])
      else
        @workspace = Workspace.select_entity_by_uri(Workspace, params[:workspace])
      end
      if @workspace.nil?
        Entity.log_debug "GridController#selectWorkspaceAndGrid " + 
                         "Invalid: can't find workspace #{params[:workspace]}"
      else
        Entity.log_debug "GridController#selectWorkspaceAndGrid workspace found name=#{@workspace.name}"
        if Entity.uuid?(params[:grid])
          @grid = Grid.select_entity_by_uuid(Grid, params[:grid])
        else
          @grid = Grid.select_entity_by_workspace_and_uri(Grid, @workspace.uuid, params[:grid])
        end
        if @grid.nil?
          Entity.log_debug "GridController#selectWorkspaceAndGrid " + 
                           "Invalid: can't find grid #{params[:grid]}"
        else
          Entity.log_debug "GridController#selectWorkspaceAndGrid grid found name=#{@grid.name}"
        end
      end
    end
  end

  def selectGridAndWorkspace
    Entity.log_debug "GridController#selectGridAndWorkspace"
    @workspace = nil
    @grid = nil
    if params[:grid].present?
      if Entity.uuid?(params[:grid])
        @grid = Grid.select_entity_by_uuid(Grid, params[:grid])
      end
      if @grid.nil?
        Entity.log_debug "GridController#selectGridAndWorkspace " + 
                         "Invalid: can't find grid #{params[:grid]}"
      else
        Entity.log_debug "GridController#selectGridAndWorkspace: grid found name=#{@grid.name}"
        @workspace = @grid.load_workspace
        if @workspace.nil?
          Entity.log_debug "GridController#selectGridAndWorkspace " + 
                           "Invalid: can't find workspace #{@grid.workspace_uuid}"
        else
          Entity.log_debug "GridController#selectGridAndWorkspace: workspace found name=#{@workspace.name}"
        end
      end
    end
  end

  def selectRow
    Entity.log_debug "GridController#selectRow"
    @row = nil
    if @grid.present? and params[:row].present?
      if Entity.uuid?(params[:row])
        @row = @row_loc = @grid.row_select_entity_by_uuid(params[:row])
        unlock_as_of_date
      else
        @row = @row_loc = @grid.row_select_entity_by_uri(params[:row])
        unlock_as_of_date
      end
      if @row.nil?
        Entity.log_debug "GridController#selectRow " +
                         "Invalid: can't find row with " +
                         "grid_uuid=#{@grid.uuid} " +
                         "and uuid=#{params[:row]}"
      else
        Entity.log_debug "GridController#selectRow: row found name=#{@row.name}"
        change_as_of_date(@row)
      end
    end
  end
end