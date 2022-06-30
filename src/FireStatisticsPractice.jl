# -----------------------------------------------------------------------------
# Разметка CSV
# -----------------------------------------------------------------------------
#        dt | Дата возгорания (ГГГГ-ММ-ДД)
# type_name | Текстовое описание
#   type_id | Код типа пожара
#       lon | Долгота возгорания
#       lat | Широта возгорания


module FireStatisticsPractice

    include("./Utils.jl")

    using .Utils

    using Base.Filesystem: ispath, mkpath
    using CSV
    using DataFrames
    using Dates

    const INPUT_PATH  = "data/"
    const OUTPUT_PATH = "output/"

    function reset_stats!( )
        global stat_total_count      = 0                        # Общее количество
        global stat_year_count       = Dict{ Int32, Int32  }( ) # По годам
        global stat_month_count      = Dict{ Int32, Int32  }( ) # По месяцам
        global stat_day_count        = Dict{ Int32, Int32  }( ) # По дням
        global stat_type_count       = Dict{ Int32, Int32  }( ) # По типу
        global stat_desc             = Dict{ Int32, String }( ) # Описание типов
        global stat_year_month_count = Dict{ Int32, Any    }( ) # По месяцам за год
        global stat_month_day_count  = Dict{ Int32, Any    }( ) # По дням за месяца
        global stat_type_year_count  = Dict{ Int32, Any    }( ) # По типу по годам
        global stat_type_month_count = Dict{ Int32, Any    }( ) # По типу по месяцам
        global stat_type_day_count   = Dict{ Int32, Any    }( ) # По типу по дням
        global stat_coordinates      = Dict{ Int32, Any    }( ) # По координатам
    end

    function get_statistics!( data )
        for row in data
            date = row.dt
            type = row.type_id
            desc = row.type_name

            year  = Dates.year( date )
            month = Dates.month( date )
            day   = Dates.day( date )

            global stat_total_count += 1

            add_key!( stat_desc, type, desc )

            update_counter!( stat_year_count,  year  )
            update_counter!( stat_month_count, month )
            update_counter!( stat_day_count,   day   )
            update_counter!( stat_type_count,  type  )

            add_key!( stat_year_month_count, year, Dict{ Int32, Int32 }( ) )
            update_counter!( stat_year_month_count[ year ], month )

            add_key!( stat_month_day_count, month, Dict{ Int32, Int32 }( ) )
            update_counter!( stat_month_day_count[ month ], day )

            add_key!( stat_type_year_count, type, Dict{ Int32, Int32 }( ) )
            update_counter!( stat_type_year_count[ type ], year )

            add_key!( stat_type_month_count, type, Dict{ Int32, Int32 }( ) )
            update_counter!( stat_type_month_count[ type ], month )

            add_key!( stat_type_day_count, type, Dict{ Int32, Int32 }( ) )
            update_counter!( stat_type_day_count[ type ], day )

            add_key!( stat_coordinates, year, Dict{ Int32, Any }( ) )
            add_key!( stat_coordinates[ year ], month, Dict{ Int32, Any }( ) )
            add_key!( stat_coordinates[ year ][ month ], day, Dict{ Int32, Any }( ) )
            add_key!( stat_coordinates[ year ][ month ][ day ], type, Dict{ Int32, Tuple{ Float32, Float32 } }( ) )

            stat_coordinates[ year ][ month ][ day ][ type ][ stat_total_count ] = ( row.lon, row.lat )
        end
    end

    function save_statistics( file_name )
        prefix = file_name * "_data_"

        begin
            stat = sort( stat_year_count )

            column_1 = [ ]
            column_2 = [ ]

            for ( key, value ) in stat
                push!( column_1, key   )
                push!( column_2, value )
            end

            CSV.write( prefix * "year.csv",
                DataFrame(
                    year  = column_1,
                    count = column_2
                )
            )
        end

        begin
            stat = sort( stat_month_count )

            column_1 = [ ]
            column_2 = [ ]

            for ( key, value ) in stat
                push!( column_1, key   )
                push!( column_2, value )
            end

            CSV.write( prefix * "month.csv",
                DataFrame(
                    month = column_1,
                    count = column_2
                )
            )
        end

        begin
            stat = sort( stat_day_count )

            column_1 = [ ]
            column_2 = [ ]

            for ( key, value ) in stat
                push!( column_1, key   )
                push!( column_2, value )
            end

            CSV.write( prefix * "day.csv",
                DataFrame(
                    day   = column_1,
                    count = column_2
                )
            )
        end

        begin
            stat = sort( stat_type_count )
            desc = sort( stat_desc )

            column_1 = [ ]
            column_2 = [ ]
            column_3 = [ ]

            for ( key, value ) in stat
                push!( column_1, key         )
                push!( column_2, desc[ key ] )
                push!( column_3, value       )
            end

            CSV.write( prefix * "type.csv",
                DataFrame(
                    type        = column_1,
                    description = column_2,
                    count       = column_3
                )
            )
        end

        begin
            stat = sort( stat_year_month_count )

            column_1 = [ ]
            column_2 = [ ]
            column_3 = [ ]

            for ( key, value ) in stat
                for ( subkey, subvalue ) in sort( value )
                    push!( column_1, key      )
                    push!( column_2, subkey   )
                    push!( column_3, subvalue )
                end
            end

            CSV.write( prefix * "year_month.csv",
                DataFrame(
                    year  = column_1,
                    month = column_2,
                    count = column_3
                )
            )
        end

        begin
            stat = sort( stat_month_day_count )

            column_1 = [ ]
            column_2 = [ ]
            column_3 = [ ]

            for ( key, value ) in stat
                for ( subkey, subvalue ) in sort( value )
                    push!( column_1, key      )
                    push!( column_2, subkey   )
                    push!( column_3, subvalue )
                end
            end

            CSV.write( prefix * "month_day.csv",
                DataFrame(
                    month = column_1,
                    day   = column_2,
                    count = column_3
                )
            )
        end

        begin
            stat = sort( stat_type_year_count )

            column_1 = [ ]
            column_2 = [ ]
            column_3 = [ ]

            for ( key, value ) in stat
                for ( subkey, subvalue ) in sort( value )
                    push!( column_1, key      )
                    push!( column_2, subkey   )
                    push!( column_3, subvalue )
                end
            end

            CSV.write( prefix * "type_year.csv",
                DataFrame(
                    type  = column_1,
                    year  = column_2,
                    count = column_3
                )
            )
        end

        begin
            stat = sort( stat_type_month_count )

            column_1 = [ ]
            column_2 = [ ]
            column_3 = [ ]

            for ( key, value ) in stat
                for ( subkey, subvalue ) in sort( value )
                    push!( column_1, key      )
                    push!( column_2, subkey   )
                    push!( column_3, subvalue )
                end
            end

            CSV.write( prefix * "type_month.csv",
                DataFrame(
                    type  = column_1,
                    month = column_2,
                    count = column_3
                )
            )
        end

        begin
            stat = sort( stat_type_day_count )

            column_1 = [ ]
            column_2 = [ ]
            column_3 = [ ]

            for ( key, value ) in stat
                for ( subkey, subvalue ) in sort( value )
                    push!( column_1, key      )
                    push!( column_2, subkey   )
                    push!( column_3, subvalue )
                end
            end

            CSV.write( prefix * "type_day.csv",
                DataFrame(
                    type  = column_1,
                    day   = column_2,
                    count = column_3
                )
            )
        end

        begin
            stat = sort( stat_coordinates )

            column_1 = [ ]
            column_2 = [ ]
            column_3 = [ ]
            column_4 = [ ]
            column_5 = [ ]
            column_6 = [ ]

            for ( year, months ) in stat
                for ( month, days ) in months
                    for ( day, types ) in days
                        for ( type, coords ) in types
                            for ( _, coord ) in coords
                                push!( column_1, year       )
                                push!( column_2, month      )
                                push!( column_3, day        )
                                push!( column_4, type       )
                                push!( column_5, coord[ 1 ] )
                                push!( column_6, coord[ 2 ] )
                            end
                        end
                    end
                end
            end

            CSV.write( prefix * "coordinates_year.csv",
                DataFrame(
                    year = column_1,
                    lon  = column_5,
                    lat  = column_6
                )
            )

            CSV.write( prefix * "coordinates_month.csv",
                DataFrame(
                    month = column_2,
                    lon   = column_5,
                    lat   = column_6
                )
            )

            CSV.write( prefix * "coordinates_year_month.csv",
                DataFrame(
                    year  = column_1,
                    month = column_2,
                    lon   = column_5,
                    lat   = column_6
                )
            )

            CSV.write( prefix * "coordinates_month_day.csv",
                DataFrame(
                    month = column_2,
                    day   = column_3,
                    lon   = column_5,
                    lat   = column_6
                )
            )

            CSV.write( prefix * "coordinates_type.csv",
                DataFrame(
                    type = column_4,
                    lon  = column_5,
                    lat  = column_6
                )
            )

            CSV.write( prefix * "coordinates_type_year.csv",
                DataFrame(
                    year = column_1,
                    type = column_4,
                    lon  = column_5,
                    lat  = column_6
                )
            )

            CSV.write( prefix * "coordinates_type_month.csv",
                DataFrame(
                    month = column_2,
                    type  = column_4,
                    lon   = column_5,
                    lat   = column_6
                )
            )
        end
    end

    function main( )
        if !ispath( OUTPUT_PATH )
            mkpath( OUTPUT_PATH )
        end

        for file in get_input( INPUT_PATH )
            reset_stats!( )

            csv = CSV.File( INPUT_PATH * file * ".csv" )

            get_statistics!( csv )

            save_statistics( OUTPUT_PATH * file )
        end
    end

    main( )

end # module
