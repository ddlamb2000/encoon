xml.instruct! :xml, :version=>"1.0"
xml.encoon do
  if @grid.present?
    if @rows.present?
      for row in @rows 
        @grid.row_export(xml, row)
      end
    end
  end
end