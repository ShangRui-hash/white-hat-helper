import React, {Component} from 'react';
import {Layout,Menu} from "antd";
import {DollarOutlined, HomeOutlined, FileOutlined} from '@ant-design/icons';
import {Link, Redirect, Route, Switch} from "react-router-dom";

import Header from "./header";
import HomePage from "./components/homePage";
import CompanyPage from "./components/companyPage";
import MissionPage from "./components/missionPage";
import AssetsPage from "./components/assetsPage";
import AssetDetail from "./components/assetDetails";

import './css/index.less'

const { Footer,Content,Sider } = Layout;

class Home extends Component {
    state = {
        collapsed: false,
    };

    //菜单折叠
    onCollapse = collapsed => {
        this.setState({ collapsed });
    };

    render() {
        const { collapsed } = this.state;
        return (
            <Layout className='home-layout'>
                <Sider collapsible collapsed={collapsed} onCollapse={this.onCollapse}>
                    <Menu theme="dark" defaultSelectedKeys={['1']} mode="inline">
                        <Menu.Item key="1" icon={<HomeOutlined />}>
                            <Link to='/home/HomePage'>
                                首页
                            </Link>
                        </Menu.Item>
                        <Menu.Item key="2" icon={<DollarOutlined />}>
                            <Link to='/home/CompanyPage'>
                                公司管理
                            </Link>
                        </Menu.Item>
                        <Menu.Item key="3" icon={<FileOutlined />}>
                            <Link to='/home/MissionPage'>
                                任务管理
                            </Link>
                        </Menu.Item>
                    </Menu>
                </Sider>
                <Layout>
                    <Header/>
                    <Content className='home-layout-content'>
                        <Switch>
                            <Route path='/home/HomePage' component={HomePage}/>
                            <Route exact path='/home/CompanyPage' component={CompanyPage}/>
                            <Route path='/home/MissionPage' component={MissionPage}/>
                            <Route exact path='/home/CompanyPage/AssetsPage' component={AssetsPage}/>
                            <Route exact path='/home/CompanyPage/AssetsPage/details' component={AssetDetail}/>
                            <Redirect to='/home/HomePage'/>
                        </Switch>
                    </Content>
                    <Footer className='home-layout-footer'>基于Nmap网络资产扫描</Footer>
                </Layout>
            </Layout>
        );
    }
}

export default Home;
