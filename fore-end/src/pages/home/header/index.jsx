import React, {Component} from 'react';
import {Button,Col,Row} from 'antd'
import './css/index.less'

class Header extends Component {
    render() {
        return (
            <Row className='home-layout-header'>
                <Col span={2} offset={22}>
                    <Button type="primary">退出</Button>
                </Col>
            </Row>
        );
    }
}

export default Header;
