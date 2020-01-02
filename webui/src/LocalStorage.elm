port module LocalStorage exposing (decodeStorage, login, logout, viewerChanges)

import Json.Decode as Decode exposing (Decoder, Value, decodeString, decodeValue, field, string)
import Json.Encode as Encode
import Viewer exposing (Viewer)


port cache : Maybe Value -> Cmd msg


port onStoreChange : (Value -> msg) -> Sub msg


decodeStorage : Value -> Maybe Viewer
decodeStorage value =
    decodeValue Decode.string value
        |> Result.andThen (decodeString storageDecoder)
        |> Result.toMaybe


login : Viewer -> Cmd msg
login viewer =
    let
        json =
            Encode.object
                [ ( "viewer"
                  , Encode.object
                        [ ( "email", Encode.string viewer.email )
                        , ( "familyName", Encode.string viewer.familyName )
                        , ( "givenName", Encode.string viewer.givenName )
                        ]
                  )
                ]
    in
    cache (Just json)


logout : Cmd msg
logout =
    cache Nothing


viewerChanges : (Maybe Viewer -> msg) -> Sub msg
viewerChanges toMsg =
    onStoreChange
        (\value ->
            toMsg (decodeFromChange value)
        )



-- INTERNAL


decodeFromChange : Value -> Maybe Viewer
decodeFromChange value =
    Decode.decodeValue
        (Decode.field "viewer" viewerDecoder)
        value
        |> Result.toMaybe


storageDecoder : Decoder Viewer
storageDecoder =
    Decode.field "viewer" viewerDecoder


viewerDecoder : Decoder Viewer
viewerDecoder =
    Decode.map3 Viewer
        (field "email" string)
        (field "familyName" string)
        (field "givenName" string)
