xml.instruct! :xml, :version=>"1.0"
xml.encoon do
  if @grid.present?
    if @row.present?
      @grid.row_export(xml, @row)
    end
    if @grid_cast.present? and @table_rows.present?
      for row in @table_rows
        @grid_cast.row_export(xml, row)
      end
    end
  end
end