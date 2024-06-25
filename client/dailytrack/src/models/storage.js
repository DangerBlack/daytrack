import * as SecureStore from 'expo-secure-store';
import { user_self_url } from './endpoints';

const API_KEY_STORAGE_KEY = 'key';
const USERNAME_STORAGE_KEY = 'username';

export async function save_key(api_key) 
{
    await SecureStore.setItemAsync(API_KEY_STORAGE_KEY, api_key);
}

export async function load_key() 
{
    let key = await SecureStore.getItemAsync(API_KEY_STORAGE_KEY);

    return key;
}

export async function save_username(username) 
{
    await SecureStore.setItemAsync(USERNAME_STORAGE_KEY, username);
}

export async function load_username(api_key) 
{
    if(!api_key)
        api_key = await load_key();
    
    let username = await SecureStore.getItemAsync(USERNAME_STORAGE_KEY);

    if(!username)
    {
        const response = await fetch(user_self_url, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Key ${api_key}`
            },
        });

        const data = await response.json();
        console.log(data);
        if(data.username)
        {
            await save_username(data.username);
            username = data.username;
        }
    }


    return username;
}