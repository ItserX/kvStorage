box.cfg {
    listen = '0.0.0.0:3301',
}

box.once('init', function()
    box.schema.space.create('json_data', {
        format = {
            {name = 'key', type = 'string'},
            {name = 'value', type = 'any'}
        }
    })
    box.space.json_data:create_index('primary',
        { type = 'TREE', parts = {'key'}})
end)