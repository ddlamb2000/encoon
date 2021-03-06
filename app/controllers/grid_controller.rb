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
  before_filter :authenticate_user!, :except => [:show, :row, :details, :list,
                                                 :attachments, :home, :credits, :history, :refresh]

  # Message used to aknowledge requests in dialogs that use an iframe. 
  OK_MSG = "<div id='ok'>OK</div>"

  # Shows a row of data as a page or as an .xml file.
  # Workspace, grid and row should be provided as parameters according to routes.
  def show
    log_debug "GridController#show"
    @filters = params[:filters]
    selectWorkspaceAndGrid(@filters)
    if @workspace.present? and @grid.present?
      @grid.load(@filters) if not @grid.loaded
      @table_columns = @grid.filtered_columns
      selectRow
      if @row.present?
        if @grid.uuid == GRID_UUID
          @grid_cast = @grid.select_grid_cast(@row.uuid)
          @attached_grids = [@grid_cast]
        else
          @attached_grids = Grid.select_referenced_grids(@grid.uuid)
        end
        @attachments = @row.attachments
        @page_title = @grid.row_title(@row)
        respond_to do |format|
          format.html do
            push_history
            render :show
            return
          end
          format.xml do
            @table_rows = @grid_cast.row_all(@filters, nil, -1, true) if @grid.uuid == GRID_UUID and @grid_cast.present?
            render :show
            return
          end
        end
      end
    end
    render :no_data, :status => 404
  end

  # Renders the content of the article through an Ajax request.
  def row
    log_debug "GridController#row"
    @filters = params[:filters]
    selectGridAndWorkspace(@filters)
    if @workspace.present? and @grid.present?
      @grid.load(@filters) if not @grid.loaded
      @table_columns = @grid.filtered_columns
      selectRow
      if @row.present?
        @attachments = @row.attachments
        render :partial => "row"
        return
      end
    end
    render :partial => "no_data", :status => 404
  end

  # Renders the details of a data row through an Ajax request.
  def details
    log_debug "GridController#details"
    selectGridAndWorkspace
    if @workspace.present? and @grid.present?
      @grid.load if not @grid.loaded
      selectRow
      if @row.present?
        @versions = @grid.row_all_versions(@row.uuid)
        @locales = @grid.row_locales(@row.uuid, @row.version, false)
        @audits = @row.all_audits
        @attachments = @row.attachments
        session[:show_details] = true
        render :partial => "details"
        return
      end
    end
    render :partial => "no_data", :status => 404
  end

  # Renders the content of a list for a given grid through an Ajax request.
  def list
    log_debug "GridController#list"
    @filters = params[:filters]
    @search = params[:search]
    @page = params[:page]
    selectGridAndWorkspace(@filters)
    if @workspace.present? and @grid.present?
      @grid.load(@filters) if not @grid.loaded 
      @table_row_count = @grid.row_count(@filters)
      @table_rows = @grid.row_all(@filters, @search, @page, true)
      @table_columns = @grid.filtered_columns
      @display_option = @grid.display_uuid
      render :partial => "list"
      return
    end
    render :partial => "no_data", :status => 404
  end

  # Renders the list of attributes for a given grid through an Ajax request.
  def attributes
    log_debug "GridController#attributes"
    @filters = params[:filters]
    @search = params[:search]
    @page = params[:page]
    selectGridAndWorkspace(@filters)
    if @workspace.present? and @grid.present?
      @grid.load(@filters) if not @grid.loaded 
      @table_row_count = @grid.row_count(@filters)
      @table_rows = @grid.row_all(@filters, @search, @page, true)
      @table_columns = @grid.filtered_columns
      @display_option = @grid.display_uuid
      render :partial => "attributes"
      return
    end
    render :partial => "no_data", :status => 404
  end

  # Renders the attachments of an article through an Ajax request.
  def attachments
    log_debug "GridController#attachments"
    selectGridAndWorkspace
    if @workspace.present? and @grid.present?
      @grid.load(@filters) if not @grid.loaded
      selectRow
      if @row.present?
        @attachments = @row.attachments
        render :partial => "attachments"
        return
      end
    end
    render :partial => "no_data", :status => 404
  end

  # Renders a creation page for a new row of data through an Ajax request.
  def new
    log_debug "GridController#new"
    @container = params[:container]
    @refresh_list = params[:refresh_list]
    @filters = params[:filters]
    selectGridAndWorkspace(@filters)
    if @workspace.present? and @grid.present?
      @grid.load(@filters) if not @grid.loaded
      @table_columns = @grid.column_all
      @row = @grid.rows.build
      @row_loc = RowLoc.new
      @grid.row_initialization(@row, @filters)
      render :partial => "edit", :locals => {:new_row => true}
      return
    end
    render :partial => "no_data", :status => 404
  end

  # Renders an edit for a given row of data through an Ajax request.
  def edit
    log_debug "GridController#edit"
    @container = params[:container]
    @refresh_list = params[:refresh_list]
    @filters = params[:filters]
    selectGridAndWorkspace(@filters)
    if @workspace.present? and @grid.present?
      @grid.load(@filters) if not @grid.loaded
      @table_columns = @grid.column_all
      selectRow
      if @row.present?
        render :partial => "edit", :locals => {:new_row => false}
        return
      end
    end
    render :partial => "no_data", :status => 404
  end

  # Creates a new row of data through an Ajax POST request.
  def create
    log_debug "GridController#create"
    saved = false
    @filters = params[:filters]
    selectGridAndWorkspace(@filters)
    if @workspace.present? and @grid.present?
      @grid.load(@filters, true)
      @row = @grid.rows.new
      @row.initialization
      if @grid.can_create_row?(@filters)
        begin
          @row.transaction do
            log_debug "GridController#create: initialize transaction"
            @row.begin = param_begin_date(@row)
            log_debug "GridController#create: populate from parameter values"
            populate_from_params
            @grid.load(@filters)
            if @grid.has_translation
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
                  @grid.create_row_loc!(@row, @row_loc, @filters)
                else
                  @grid.row_validate(@row, Grid::PHASE_CREATE, @filters)
                  log_debug "GridController#create: rollback!"
                  raise ActiveRecord::Rollback
                end
              end
            end
            log_debug "GridController#create: row_validate"
            if @grid.row_validate(@row, Grid::PHASE_CREATE, @filters)
              log_debug "GridController#create: create_row!"
              @grid.create_row!(@row, @filters)
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
      else
        log_security_warning "GridController#create can_create_row? == false"
        @row.errors.add(:uuid, I18n.t('error.cant_create'))
      end
      if saved
        list
        return
      else
        log_debug "GridController#create: error, @row.errors=#{@row.errors.inspect}"
        render :json => @row.errors, :status => :unprocessable_entity
        return
      end
    end
    render :partial => "no_data", :status => 404
  end

  # Updates a given row of data through an Ajax POST request.
  def update
    log_debug "GridController#update"
    saved = false
    @refresh_list = params[:refresh_list]
    @filters = params[:filters]
    @container = params[:container]
    selectGridAndWorkspace(@filters)
    if @workspace.present? and @grid.present?
      @grid.load(@filters) if not @grid.loaded
      selectRow
      if @row.present?
        if @grid.can_update_row?(@row)
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
                  @row.begin = param_begin_date(@row)
                  if multiple_versions_ok?
                    @row.enabled = params[:enabled]
                    log_debug "GridController#update populate from parameter values"
                    populate_from_params
                    log_debug "GridController#update row_validate"
                    if @grid.row_validate(@row, Grid::PHASE_NEW_VERSION, @filters)
                      log_debug "GridController#update create_row!"
                      @grid.create_row!(@row, @filters)
                      if @grid.has_translation
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
                          if @grid.row_loc_validate(@row, @row_loc, Grid::PHASE_NEW_VERSION)
                            log_debug "GridController#update create_row_loc!"
                            @grid.create_row_loc!(@row, @row_loc, @filters)
                          else
                            log_debug "GridController#update row_validate"
                            @grid.row_validate(@row, Grid::PHASE_UPDATE, @filters)
                            log_debug "GridController#update: rollback!"
                            raise ActiveRecord::Rollback
                          end
                        end
                      end
                      log_debug "GridController#update row_update_dates!"
                      @grid.row_update_dates!(@row.uuid)
                      saved = true
                    else
                      raise ActiveRecord::Rollback
                    end
                  end
                else
                  log_debug "GridController#update update existing version"
                  @row.begin = param_begin_date(@row)
                  if multiple_versions_ok?
                    log_debug "GridController#update versions OK"
                    @row.enabled = params[:enabled]
                    log_debug "GridController#update populate from parameter values"
                    populate_from_params
                    log_debug "GridController#update row_validate"
                    if @grid.row_validate(@row, Grid::PHASE_UPDATE, @filters)
                      log_debug "GridController#update update_row!"
                      @grid.update_row!(@row)
                      if @grid.has_translation
                        for @row_loc in @grid.row_loc_select_entity_by_uuid(@row.uuid, version)
                          log_debug "GridController#update locale=#{@row_loc.locale}"
                          if @row_loc.locale == I18n.locale.to_s or 
                             @row_loc.base_locale == I18n.locale.to_s
                            log_debug "GridController#update update row_loc values"
                            @row_loc.base_locale = I18n.locale.to_s
                            @row_loc.name = params[:name]
                            @row_loc.description = params[:description]
                            log_debug "GridController#update row_loc_validate"
                            if @grid.row_loc_validate(@row, @row_loc, Grid::PHASE_UPDATE)
                              log_debug "GridController#update update_row_loc!"
                              @grid.update_row_loc!(@row, @row_loc)
                            else
                              log_debug "GridController#update: rollback!(3)"
                              raise ActiveRecord::Rollback
                            end
                          end
                        end
                      end
                      log_debug "GridController#update row_update_dates!"
                      @grid.row_update_dates!(@row.uuid)
                      saved = true
                    else
                      log_debug "GridController#update: rollback!(4)"
                      raise ActiveRecord::Rollback
                    end
                  end
                end
              end
            rescue ActiveRecord::RecordInvalid => invalid
              log_debug "GridController#update: invalid=#{invalid.inspect}"
              saved = false
            rescue Exception => invalid
              log_error "GridController#update", invalid
              saved = false
            end
          end
        else
          log_security_warning "GridController#update can_update_data? == false"
          @row.errors.add(:uuid, I18n.t('error.cant_update'))
        end
        if saved
          log_debug "GridController#update saved"
          change_as_of_date(@row)
          if @container == ""
            log_debug "GridController#update no refresh"
            render :nothing => true
            return
          else
            log_debug "GridController#update refresh"
            @refresh_list ? list : row
            return
          end
        else
          log_debug "GridController#update: error, params=#{params.inspect}"
          render :json => @row.errors, :status => :unprocessable_entity
          return
        end
      end
    end
    render :partial => "no_data", :status => 404
  end

  # Renders an upload page for a new attachment to the selected row of data through an Ajax request.
  def attach
    log_debug "GridController#attach"
    selectGridAndWorkspace
    if @workspace.present? and @grid.present?
      @grid.load if not @grid.loaded
      selectRow
      if @row.present?
        @attachment = @row.attachments.new
        render :partial => "attach"
        return
      end
    end
    render :partial => "no_data", :status => 404
  end

  # Attaches a document to the select data row.
  # Attachments are managed using PaperClip, through attribute 'document'.
  def save_attachment
    log_debug "GridController#save_attachment"
    saved = false
    selectGridAndWorkspace
    if @workspace.present? and @grid.present?
      @grid.load if not @grid.loaded
      selectRow
      if @row.present?
        @attachment = @row.attachments.new
        if @grid.can_update_row?(@row)
          begin
            if params[:document].blank?
              @attachment.errors.add(:document_file_name, I18n.t('error.required', :column => I18n.t('field.file')))
            else
              @row.transaction do
                @row.remove_attachment!(params[:document])
                @attachment.original_file_name = (params[:document]).original_filename
                @attachment.document = params[:document]
                @attachment.create_user_uuid = Entity.session_user_uuid
                @attachment.save!
                @grid.update_row!(@row, Audit::ATTACH)
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
        else
          log_security_warning "GridController#save_attachment: can_update_data? == false"
          @attachment.errors.add(:uuid, I18n.t('error.cant_update'))
        end
        if saved
          log_debug "GridController#save_attachment: saved"
          render :text => OK_MSG
          return
        else
          log_debug "GridController#save_attachment: error"
          render :json => @attachment.errors, :status => :unprocessable_entity
          return
        end
      end
    end
    render :partial => "no_data", :status => 404
  end

  # Removes attached document to the select data row.
  def delete_attachment
    log_debug "GridController#delete_attachment"
    saved = false
    selectGridAndWorkspace
    if @workspace.present? and @grid.present?
      @grid.load if not @grid.loaded
      selectRow
      if @row.present?
        begin
          @row.transaction do
            @attachment = @row.attachments.find(params[:id])
            if @attachment.present?
              if @grid.can_update_row?(@row)
                @attachment.delete
                @grid.update_row!(@row, Audit::DETACH)
                saved = true
              else
                log_security_warning "GridController#delete_attachment can_update_data? == false"
                @attachment.errors.add(:uuid, I18n.t('error.cant_update'))
              end
            end
          end
        rescue ActiveRecord::RecordInvalid => invalid
          log_debug "GridController#delete_attachment: invalid=#{invalid.inspect}"
          saved = false
        rescue Exception => invalid
          log_error "GridController#delete_attachment", invalid
          saved = false
        end
        if saved
          row
          return
        end
      end
    end
    render :partial => "no_data", :status => 404
  end

  # Renders the import dialog page through an Ajax request.
  def import
    log_debug "GridController#import"
    selectGridAndWorkspace
    if @workspace.present? and @grid.present?
      @grid.load if not @grid.loaded
      render :partial => "import"
      return
    end
    render :partial => "no_data", :status => 404
  end
  
  # Renders the upload dialog page through an Ajax request.
  def upload
    log_debug "GridController#upload"
    saved = false
    selectGridAndWorkspace
    if @workspace.present? and @grid.present?
      @grid.load if not @grid.loaded
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
      if saved
        log_debug "GridController#upload: saved"
        flash[:notice] = I18n.t('message.uploaded',
                                :record_count => @upload.records,
                                :insert_count => @upload.inserted,
                                :update_count => @upload.updated)
        render :text => OK_MSG
        return
      else
        log_debug "GridController#upload: error, @upload.errors=#{@upload.errors.inspect}"
        render :json => @upload.errors, :status => :unprocessable_entity
        return
      end
    end
    render :partial => "no_data", :status => 404
  end

  # Renders the home page using hard-coded references.
  def home
    params[:workspace] = SYSTEM_WORKSPACE_URI
    params[:grid] = HOME_PAGE_UUID
    params[:row] = HOME_WELCOME_UUID
    show
  end

  # Renders credits page using hard-coded references.
  def credits
    params[:workspace] = SYSTEM_WORKSPACE_URI
    params[:grid] = HOME_PAGE_UUID
    params[:row] = HOME_CREDITS_UUID
    show
  end

  # Renders history of navigation through an Ajax request.
  def history
    session[:show_history] = true
    render :partial => "history"
  end

  # Sets the session flag referenced by the given parameter through an Ajax request.
  def set
    session[params[:flag]] = true unless params[:flag].nil?
    render :nothing => true
  end

  # Unsets the session flag referenced by the given parameter through an Ajax request.
  def unset
    session[params[:flag]] = false unless params[:flag].nil?
    render :nothing => true
  end

  # Forces page to be refreshed after an as of date or locale change.
  def refresh
    log_debug "GridController#refresh"
    redirect_to session[:last_url]
  end

private

  # Returns the begin date is present as part of parameters.
  # Returns a default value when the begin date is absent from the parameters.
  def param_begin_date(row)
    begin
      params[:begin].present? ? 
        Date.strptime(params[:begin], t('date.formats.default')) :
        Entity.begin_of_time
    rescue
      @row.errors.add(:begin, I18n.t('error.badformat', :column => I18n.t('field.begin_date')))
      Entity.begin_of_time
    end
  end

  # Sets parameter values into data row using physical column names.
  # Used for creation and update of rows. 
  def populate_from_params
    if @grid.present? and @row.present?
      for column in @grid.column_all
        value = params["row_#{column.physical_column}"]
        log_debug "GridController#populate_from_params #{column.physical_column} = #{value}"
        @row.write_value(column, value)
        log_debug "GridController#populate_from_params verify value : #{@row.read_value(column)}"
      end
      @row.uri = params["uri"]
      @row.clean_uri!
    end
  end

  # Controls if version of the row exists or not in the database in order to avoid duplicates.
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

  # Selects workspace, then grid information based on paramaters.
  # The workspace is selected based on given uri or uuid.
  # The grid is selected for the given workspace based on uri or uuid.
  def selectWorkspaceAndGrid(filters=nil)
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
          params[:workspace] = SYSTEM_WORKSPACE_URI
          params[:grid] = WORKSPACE_UUID
        elsif not params[:grid].nil? and params[:row].nil?
          Entity.log_debug "GridController#selectWorkspaceAndGrid grid only"
          params[:row] = params[:grid]
          params[:workspace] = SYSTEM_WORKSPACE_URI
          params[:grid] = GRID_UUID
        end
      end 
      Entity.log_debug "GridController#selectWorkspaceAndGrid select workspace"
      if Entity.uuid?(params[:workspace])
        @workspace = Workspace.select_entity_by_uuid(Workspace, params[:workspace])
      else
        @workspace = Workspace.select_entity_by_uri(Workspace, params[:workspace])
      end
      if @workspace.nil?
        Entity.log_debug "GridController#selectWorkspaceAndGrid can't find workspace #{params[:workspace]}"
      else
        Entity.log_debug "GridController#selectWorkspaceAndGrid workspace found name=#{@workspace.name}"
        Entity.log_debug "GridController#selectWorkspaceAndGrid select grid"
        if Entity.uuid?(params[:grid])
          @grid = Grid.select_entity_by_uuid(Grid, params[:grid], filters)
        else
          @grid = Grid.select_entity_by_workspace_and_uri(Grid, @workspace.uuid, params[:grid], filters)
        end
        if @grid.nil?
          Entity.log_debug "GridController#selectWorkspaceAndGrid can't find grid #{params[:grid]}"
        else
          Entity.log_debug "GridController#selectWorkspaceAndGrid grid found name=#{@grid.name}"
        end
      end
    end
  end

  # Selects grid, then workspace information based on paramaters.
  # The grid is selected based on uuid only.
  # The workspace is selected based on grid information.
  def selectGridAndWorkspace(filters = nil)
    Entity.log_debug "GridController#selectGridAndWorkspace"
    @workspace = nil
    @grid = nil
    if params[:grid].present?
      Entity.log_debug "GridController#selectGridAndWorkspace select grid"
      if Entity.uuid?(params[:grid])
        @grid = Grid.select_entity_by_uuid(Grid, params[:grid], filters)
      end
      if @grid.nil?
        Entity.log_debug "GridController#selectGridAndWorkspace can't find grid #{params[:grid]}"
      else
        Entity.log_debug "GridController#selectGridAndWorkspace: grid found name=#{@grid.name}"
        @grid.load(filters) if not @grid.loaded
        @workspace = @grid.workspace
        if @workspace.nil?
          Entity.log_debug "GridController#selectGridAndWorkspace can't find workspace #{@grid.workspace_uuid}"
        else
          Entity.log_debug "GridController#selectGridAndWorkspace: workspace found name=#{@workspace.name}"
        end
      end
    end
  end

  # Selects row information based on paramaters.
  # Row is selected for a given grid.
  def selectRow
    Entity.log_debug "GridController#selectRow"
    @row = nil
    if @grid.present? and params[:row].present?
      if Entity.uuid?(params[:row])
        @row = @row_loc = @grid.row_select_entity_by_uuid(params[:row])
      else
        @row = @row_loc = @grid.row_select_entity_by_uri(params[:row])
      end
      if @row.nil?
        Entity.log_debug "GridController#selectRow can't find row #{params[:row]}"
      else
        Entity.log_debug "GridController#selectRow: row found name=#{@row.name}"
        if @grid.uuid == GRID_UUID
          Entity.log_debug "GridController#selectRow: reselect workspace"
          @workspace = Workspace.select_entity_by_uuid(Workspace, @row.workspace_uuid)
        end
        change_as_of_date(@row)
      end
    end
  end
end