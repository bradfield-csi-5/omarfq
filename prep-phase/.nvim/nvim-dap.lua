local dap = require('dap');

dap.adapters.codelldb = {
    type = 'server',
    port = "13000",
    host = "127.0.0.1"
}

dap.configurations.c = {
    {
        name = "Launch file",
        type = "codelldb",
        request = "launch",
        program = function()
            return vim.fn.input('Path to executable: ', vim.fn.getcwd() .. '/', 'file')
        end,
        cwd = '${workspaceFolder}',
        stopOnEntry = false,
    },
}
