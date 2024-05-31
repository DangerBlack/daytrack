import React from 'react';
import { View, Text, TextInput, TouchableOpacity, ActivityIndicator } from 'react-native';
import styles from '../styles/style';
import { useState } from 'react';
import { magic_link_url } from '../models/endpoints';
const LoginForm = () => {
    const [email, setEmail] = useState('');
    const [is_loading, set_is_loading] = useState(false);

    const handleLogin = async () => {
        set_is_loading(true);
        try
        {
            console.log('Login button pressed!');
            const response = await fetch(magic_link_url, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ email })
            });
            const data = await response.json();
            console.log(data);
        }
        catch (error)
        {
            console.error(error);
        }
        finally
        {
            set_is_loading(false);
        }   
    };

    return (
        <View style={styles.container}>
            <View style={styles.formContainer}>
                <Text style={styles.title}>daily track</Text>
                <TextInput
                    style={styles.input}
                    placeholder="email"
                    keyboardType="email-address"
                    autoCapitalize="none"
                    autoCorrect={false}
                    onChangeText={text => setEmail(text)}
                />
                {
                    is_loading ? 
                    <ActivityIndicator size="large" color="#007BFF" /> :
                    <TouchableOpacity style={styles.button} onPress={handleLogin} >
                        <Text style={styles.buttonText}>Login</Text>
                    </TouchableOpacity>
                }
            </View>
        </View>
    );
};

export default LoginForm;