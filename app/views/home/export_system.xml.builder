xml.instruct! :xml, :version=>"1.0"
xml.encoon do
  @workspace.export(xml)
  Grid.all(@workspace.grids, session[:as_of_date]).each do |grid| 
    grid.load_cached_grid_structure 
    grid.export(xml)
    grid.mapping_all.each do |mapping|
      grid.mapping_export(xml, mapping)
    end
    grid.column_all.each do |column| 
      grid.column_export(xml, column)
      column.column_mapping_all.each do |mapping|
        grid.column_mapping_export(xml, column, mapping)
      end
    end
    if grid.uuid == Column::ROOT_DATA_KIND_UUID or 
       grid.uuid == Column::ROOT_GRID_DISPLAY_OPTION_UUID or
       grid.uuid == Grid::ROOT_COUNTRY_UUID or
       grid.uuid == Role::ROOT_UUID
      grid.row_all(nil, nil, -1).each do |row|
        grid.row_export(xml, row)
      end
    end
  end
end