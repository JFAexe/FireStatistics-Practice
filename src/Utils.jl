module Utils

    using Base.Filesystem: readdir

    export get_input, add_key!, update_counter!

    function get_input( path )
        map( ( file ) -> replace( file, ".csv" => "" ), readdir( path ) )
    end

    function add_key!( dictionary, key, value = 0 )
        if haskey( dictionary, key )
            return
        end

        dictionary[ key ] = value
    end

    function update_counter!( dictionary, key, value = 1 )
        add_key!( dictionary, key )

        dictionary[ key ] += value
    end

end # module
