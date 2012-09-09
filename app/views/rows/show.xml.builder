xml.instruct! :xml, :version=>"1.0"
xml.encoon do
  if @grid.present?
    if @row.present?
      @grid.row_export(xml, @row)
    else
      for row in @rows 
        @grid.row_export(xml, row)
      end
    end
  end
end