import React, {Component} from 'react';
import {Input,Col} from "antd";
import './css/index.less'
const {Search} = Input

class HomePage extends Component {
    render() {
        return (
            <div>
                <Col className='home-layout-main' offset={7} span={10}>
                    <h1>网络资产扫描</h1>
                    <Search className='main-search' size={'large'}/>
                </Col>
            </div>
        );
    }
}

export default HomePage;
