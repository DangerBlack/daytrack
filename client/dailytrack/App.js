import React, { useMemo, useState, useEffect } from 'react';
import * as Linking from 'expo-linking';
import { createDrawerNavigator } from '@react-navigation/drawer';
import { NavigationContainer } from '@react-navigation/native';
import LoginForm from './src/views/login';
import HomePage from './src/views/home';
import { save_key, load_key, load_username } from './src/models/storage';

const Drawer = createDrawerNavigator();

async function check_if_logged_in(set_is_logged) 
{
  const api_key = await load_key();
  const username = await load_username(api_key);

  if(api_key)
  {
    console.log("use is logged with:", api_key)
    set_is_logged(true);
  }

  if(!username)
  {
    console.log("no username")
  }
}

export default function App() {
  const [is_logged, set_is_logged] = useState(false);
  const url = Linking.useURL();
  console.log(url);

  useMemo(async () => 
  {
    try
    {
      if(url)
      {
        const api_key = url.split('key=')[1];
        if(api_key)
        {
          console.log("storing secret:", api_key)
          await save_key(api_key);
        }
      }
    }
    catch(e)
    {
      console.error(e);
    }
  }, [url]);

  check_if_logged_in(set_is_logged).catch(console.error);

  return (
    <NavigationContainer>
      <Drawer.Navigator initialRouteName="Login">
          {is_logged || <Drawer.Screen name="Login" component={LoginForm} />}
          {!is_logged || <Drawer.Screen name="Home" component={HomePage} />}
      </Drawer.Navigator>
    </NavigationContainer>
  );
}

