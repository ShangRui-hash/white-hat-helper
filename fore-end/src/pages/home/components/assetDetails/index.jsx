import React, {Component} from 'react';
import {Card,Row,Col} from "antd";

import './css/index.less'

class AssetDetail extends Component {
    render() {
        return (
            <Card
                className='data-card'
                title={<p className='title-font'>资产详情</p>}
                bordered={false}
            >
                <Row gutter={[15,15]}>
                    <Col span={12}>
                        <Card
                            className='details-card'
                            title={<p className='title-font'>基本信息</p>}
                            size={'small'}
                        >
                            ?????<br/>
                            ????<br/>
                            ????<br/>
                            ???<br/>
                        </Card>
                    </Col>
                    <Col span={12}>
                        <Card
                            className='details-card'
                            title={<p className='title-font'>开放的端口</p>}
                            size={'small'}
                        >
                            ?????<br/>
                            ????<br/>
                            ????<br/>
                            ???<br/>
                        </Card>
                    </Col>
                    <Col span={12}>
                        <Card
                            className='details-card'
                            title={<p className='title-font'>漏洞</p>}
                            size={'small'}
                        >
                            ?????<br/>
                            ????<br/>
                            ????<br/>
                            ???<br/>
                        </Card>
                    </Col>
                    <Col span={12}>
                        <Card
                            className='details-card'
                            title={<p className='title-font'>基本信息</p>}
                            size={'small'}
                        >
                            ?????<br/>
                            ????<br/>
                            ????<br/>
                            ???<br/>
                        </Card>
                    </Col>
                </Row>
            </Card>
        );
    }
}

export default AssetDetail;
