import React, { useState, useEffect,useCallback } from 'react';
import { View, TouchableOpacity, FlatList, Text, ActivityIndicator, RefreshControl } from 'react-native';
// import { Icon } from 'react-native-elements';
import styles from '../styles/style';
import { track_url, event_url } from '../models/endpoints';
import { load_key, load_username } from '../models/storage';
import { createStackNavigator } from '@react-navigation/stack';
import CreateNewTrack from './create_track';

const HomeStack = createStackNavigator();

const ButtonGrid = ({ buttons, onButtonPress, onButtonLongPress, is_loading, set_is_loading, set_tracks }) => {
    const renderButton = ({ item }) => (
        <TouchableOpacity style={styles.buttonInGrid} onPress={() => onButtonPress(item.id, item.name)} onLongPress={() => onButtonLongPress(item.id)}>
          <Text style={styles.buttonText}>{item.name}</Text>
        </TouchableOpacity>
    );

    const onRefresh = useCallback(() => {
        load_tracks(set_is_loading, set_tracks).catch(console.error);
      }, [buttons]);

    return (
        <FlatList
        data={buttons}
        renderItem={renderButton}
        keyExtractor={item => item.id}
        numColumns={2}
        contentContainerStyle={styles.buttonGrid}
        refreshControl={
            <RefreshControl refreshing={is_loading} onRefresh={onRefresh} />
          }
        />
    );
};

async function load_tracks(set_is_loading, set_buttons) 
{
    set_is_loading(true);
    try {
        const api_key = await load_key();
        const response = await fetch(track_url, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Key ${api_key}`
            },
        });
        const data = await response.json();
        console.log(data);
        if(data.items)
            set_buttons(data.items);
    } catch (error) {
      console.error(error);
    } finally {
      set_is_loading(false);
    }
}

async function track_event(id, name) 
{
    try {
        const api_key = await load_key();
        const username = await load_username(api_key);
        const response = await fetch(`${event_url}/${username}/${name}?key=${api_key}`, {
            method: 'POST',
        });
        const data = await response.json();
        console.log(data);
    } catch (error) {
      console.error(error);
    } 
}

const HomePage = ({ navigation }) => {
    const [is_loading, set_is_loading] = useState(false);
    const [tracks, set_tracks] = useState([
    ]);
    
    const handleButtonPress = async (id, name) => {
        console.log(`Button ${id} - ${name} pressed`);
        track_event(id, name).catch(console.error);
    };
    
    const handleButtonLongPress = (id) => {
        console.log(`Button ${id} loong pressed`);
    };

    const handleAddTrack = () => {
        navigation.navigate('CreateNewTrack');
    };

    useEffect(() => 
    {
        load_tracks(set_is_loading, set_tracks).catch(console.error);
        return
    }, []);
    
    return (
    <View style={styles.homeContainer}>
        <ButtonGrid 
            buttons={tracks} 
            onButtonPress={handleButtonPress} 
            onButtonLongPress={handleButtonLongPress}
            is_loading={is_loading} 
            set_is_loading={set_is_loading} 
            set_tracks={set_tracks} 
        />
            <TouchableOpacity style={styles.fab} onPress={handleAddTrack}>
            <Text style={styles.buttonText}>+</Text>
        </TouchableOpacity>
    </View>
    );
};

const HomePageNav = ({ navigation }) => {
    return (
        <HomeStack.Navigator>
            <HomeStack.Screen
                name="HomePage"
                component={HomePage}
                options={{ headerShown: false }}
            />
            <HomeStack.Screen 
                name="CreateNewTrack" 
                component={CreateNewTrack} 
                options={{ title: 'Create new track' }}
            />
        </HomeStack.Navigator>
    );
};

export default HomePageNav;