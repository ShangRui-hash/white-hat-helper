/*资产页面*/
import React, {Component} from 'react';
import {Card} from "antd";

class PropertyPage extends Component {
    render() {
        return (
            <Card
                className='data-card'
                title={<p className='title-font'>公司资产</p>}
                bordered={false}
            >
                资产
            </Card>
        );
    }
}

export default PropertyPage;
