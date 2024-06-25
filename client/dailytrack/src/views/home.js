import React, { useState, useEffect } from 'react';
import { View, TouchableOpacity, FlatList, Text, ActivityIndicator } from 'react-native';
// import { Icon } from 'react-native-elements';
import styles from '../styles/style';
import { track_url } from '../models/endpoints';
import { load_key } from '../models/storage';


const ButtonGrid = ({ buttons, onButtonPress }) => {
    const renderButton = ({ item }) => (
        <TouchableOpacity style={styles.buttonInGrid} onPress={() => onButtonPress(item.id)}>
          <Text style={styles.buttonText}>{item.name}</Text>
        </TouchableOpacity>
    );

    return (
        <FlatList
        data={buttons}
        renderItem={renderButton}
        keyExtractor={item => item.id}
        numColumns={2}
        contentContainerStyle={styles.buttonGrid}
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

async function add_track(set_is_loading, set_buttons) 
{
    console.log('Adding a track')
    set_is_loading(true);
    try {
        const api_key = await load_key();
        const response = await fetch(track_url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Key ${api_key}`
            },
            body: JSON.stringify({ name: 'lol' }),
        });
        console.log('Track added')
        const data = await response.json();
        console.log(data);

    } catch (error) {
      console.error(error);
    } finally {
      set_is_loading(false);
    }
}

const HomePage = () => {
    const [is_loading, set_is_loading] = useState(false);
    const [tracks, set_tracks] = useState([
    ]);
    
    const handleButtonPress = (id) => {
        console.log(`Button ${id} pressed`);
    };

    const handleAddTrack = () => {
        add_track(set_is_loading, set_tracks).catch(console.error);
    };

    useEffect(() => 
    {
        load_tracks(set_is_loading, set_tracks).catch(console.error);
        return
    }, []);
    
      return (
        <View style={styles.homeContainer}>
          {
            is_loading ?
            <ActivityIndicator size="large" color="#007BFF" />
            : <>
            <ButtonGrid buttons={tracks} onButtonPress={handleButtonPress} />
                <TouchableOpacity style={styles.fab} onPress={handleAddTrack}>
                <Text style={styles.buttonText}>+</Text>
            </TouchableOpacity>
            </>
        }
        </View>
      );
};

export default HomePage;