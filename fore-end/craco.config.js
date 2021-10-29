const CracoLessPlugin = require('craco-less');

module.exports = {
    plugins: [
        {
            plugin: CracoLessPlugin,
            options: {
                lessLoaderOptions: {
                    lessOptions: {
                        modifyVars: {
                            '@primary-color': '#1890ff' ,
                            '@border-color-base': '#05f2f2'
                        },
                        javascriptEnabled: true,
                    },
                },
            },
        },
    ],
};
