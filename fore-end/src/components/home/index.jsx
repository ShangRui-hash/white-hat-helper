import React, {Component} from 'react';
import {Layout} from "antd";
import './css/index.less'

const { Header, Footer,Content } = Layout;

class Home extends Component {
    render() {
        return (
            <Layout className='home-layout'>
                <Header>Header</Header>
                <Content>Content</Content>
                <Footer>Footer</Footer>
            </Layout>
        );
    }
}

export default Home;
