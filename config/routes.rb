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
Encoon::Application.routes.draw do
  root :to => 'grid#home'
  devise_for :users
  match '*flag/__set' => 'grid#set', :as => 'set'
  match '*flag/__unset' => 'grid#unset', :as => 'unset'
  match '/*grid/__list' => 'grid#list', :as => 'list'
  match '/*grid/__new' => 'grid#new', :as => 'new'
  match '/*grid/__create' => 'grid#create', :as => 'create', :via => [:post]
  match '/*grid/__import' => 'grid#import', :as => 'import'
  match '/*grid/__upload' => 'grid#upload', :as => 'upload', :via => [:post]
  match '/*grid/__attributes' => 'grid#attributes', :as => 'attributes'
  match '/*grid/*row/__details' => 'grid#details', :as => 'details'
  match '/*grid/*row/__edit' => 'grid#edit', :as => 'edit'
  match '/*grid/*row/__update' => 'grid#update', :as => 'update', :via => [:post]
  match '/*grid/*row/__attach_document' => 'grid#attach_document', :as => 'attach_document'
  match '/*grid/*row/__save_attachment' => 'grid#save_attachment', :as => 'save_attachment'
  match '/*grid/*row/__delete_attachment' => 'grid#delete_attachment', :as => 'delete_attachment'
  match '/*workspace' => 'grid#show', :as => 'show_workspace'
  match '/*workspace/*grid' => 'grid#show', :as => 'show_grid'
  match '/*workspace/*grid/*row.xml' => 'grid#show', :format => :xml, :as => 'export_row_xml'
  match '/*workspace/*grid/*row' => 'grid#show', :as => 'show'
  match '/*workspace/*grid.xml' => 'grid#show', :format => :xml, :as => 'export_list_xml'
  match '/__history' => 'grid#history', :as => 'history'
  match '/__refresh' => 'grid#refresh', :as => 'refresh'
end