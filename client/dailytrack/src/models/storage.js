import * as SecureStore from 'expo-secure-store';

const API_KEY_STORAGE_KEY = 'key';

export async function save_key(api_key) 
{
    await SecureStore.setItemAsync(API_KEY_STORAGE_KEY, api_key);
}

export async function load_key() 
{
    let key = await SecureStore.getItemAsync(API_KEY_STORAGE_KEY);

    return key;
}