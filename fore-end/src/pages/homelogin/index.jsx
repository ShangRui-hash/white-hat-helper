import React, {Component} from 'react';
import {Layout,Form,Input,Col,Avatar,Button} from "antd";
import { UserOutlined, LockOutlined } from '@ant-design/icons';

import Header from './header/index'
import './css/index.less'

const { Footer, Content } = Layout;
const {Password} = Input

class HomeLogin extends Component {
    render() {
        return (
            <Layout className='login-layout'>
                <Header/>
                <Content className='login-layout-content'>
                    <Col className = 'login-form' offset={8} span={8}>
                        <div>
                            <Avatar size={100} style={{margin:'10px'}} src="https://joeschmoe.io/api/v1/random"/>
                            <Form style={{margin:'10px auto',width:'300px'}}>
                                <Form.Item name='username'>
                                    <Input prefix={<UserOutlined/>} placeholder="用户名" allowClear/>
                                </Form.Item>
                                <Form.Item name='password'>
                                    <Password prefix={<LockOutlined/>} placeholder="密码" allowClear/>
                                </Form.Item>
                                <Form.Item>
                                    <Button type='primary' block>登录</Button>
                                </Form.Item>
                            </Form>
                        </div>
                    </Col>
                </Content>
                <Footer className='login-layout-footer'>基于Nmap网络资产扫描</Footer>
            </Layout>
        );
    }
}

export default HomeLogin;
