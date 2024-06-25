import { StyleSheet } from 'react-native';

const primaryColor = '#007BFF';
const secondaryColor = '#333';
const accentColor = '#FFC107';

const smallTextSize = 16;
const defaultTextSize = 18;
const titleTestSize = 24;

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: 'white',
        justifyContent: 'center',
        alignItems: 'center',
    },
    formContainer: {
        width: '80%',
    },
    title: {
        textAlign: 'center',
        fontSize: titleTestSize,
        fontWeight: 'bold',
        marginBottom: 40,
        fontFamily: 'Roboto',
    },
    label: {
        fontSize: smallTextSize,
        marginBottom: 5,
        color: secondaryColor,
    },
    input: {
        borderWidth: 1,
        borderColor: '#ccc',
        borderRadius: 5,
        paddingVertical: 10,
        paddingHorizontal: 15,
        marginBottom: 10,
    },
    button: {
        backgroundColor: '#007BFF',
        paddingVertical: 12,
        borderRadius: 5,
    },
    buttonText: {
        textAlign: 'center',
        color: 'white',
        fontWeight: 'bold',
        fontSize: defaultTextSize,
    },

    homeContainer: {
        flex: 1,
        backgroundColor: 'white',
        justifyContent: 'center',
        alignItems: 'center',
    },
    buttonGrid: {
        justifyContent: 'center',
        alignItems: 'center',
    },
    buttonInGrid: {
        backgroundColor: '#007BFF',
        width: 100,
        height: 100,
        margin: 10,
        borderRadius: 10,
        overflow: 'hidden',
        justifyContent: 'center',
        alignItems: 'center',
    },
    buttonImage: {
        width: '100%',
        height: '100%',
    },
    fab: {
        position: 'absolute',
        bottom: 20,
        right: 20,
        width: 60,
        height: 60,
        borderRadius: 30,
        backgroundColor: '#6200EE',
        justifyContent: 'center',
        alignItems: 'center',
        elevation: 5,
    },
});

export default styles;