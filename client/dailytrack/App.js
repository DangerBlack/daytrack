import React from 'react';
import { createDrawerNavigator } from '@react-navigation/drawer';
import { NavigationContainer } from '@react-navigation/native';
import LoginForm from './src/views/login';
import HomePage from './src/views/home';

const Drawer = createDrawerNavigator();

const LoginPage = () => {
    return (
        <NavigationContainer>
            <Drawer.Navigator initialRouteName="Login">
                <Drawer.Screen name="Login" component={LoginForm} />
                <Drawer.Screen name="Home" component={HomePage} />
            </Drawer.Navigator>
        </NavigationContainer>
    );
};

export default function App() {
  return (
    <LoginPage />
  );
}

