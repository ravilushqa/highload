box.cfg {
    listen = 3302;
    -- 1 – SYSERROR
    -- 2 – ERROR
    -- 3 – CRITICAL
    -- 4 – WARNING
    -- 5 – INFO
    -- 6 – VERBOSE
    -- 7 – DEBUG
    log_level = 2;
}

box.schema.user.grant('guest', 'read,write,execute', 'universe')

local function bootstrap()

    if not box.space.mysqldaemon then
        s = box.schema.space.create('mysqldaemon')
        s:create_index('primary',
                { type = 'tree', parts = { 1, 'unsigned' }, if_not_exists = true })
    end

    if not box.space.mysqldata then
        t = box.schema.space.create('mysqldata')
        t:create_index('primary',
                { type = 'tree', parts = { 1, 'unsigned' }, if_not_exists = true })
    end

end

bootstrap()