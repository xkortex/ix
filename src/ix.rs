mod slice;

extern crate clap;
use clap::{Arg, App, SubCommand};
use slice::FlexibleSlice;
use std::convert::TryFrom;

fn main() {
    let matches = App::new("ix")
        .version("0.1")
        .author("Mike M")
        .about("A better slicing tool")
        .args_from_usage(
            "-c, --config=[FILE]    'Sets a custom config file'
            <SLICE>                 'Slice notation'
            [FILES]...                  'File to process (otherwise read stdin)'
            -v...                   'Sets the level of verbosity'")
        .get_matches();

    println!("{:?}", matches);
    let slice = matches.value_of("SLICE").unwrap();
    println!("{:?}", slice);
    let flex_slice = FlexibleSlice::from(String::from(slice));
    println!("{:?}", flex_slice);
    // let hmm = FlexibleSlice::try_from(":");
    println!("{:?}", slice::FlexibleSlice::try_from(":").unwrap());
    println!("{:?}", slice::FlexibleSlice::try_from("1:2:-3").unwrap());
    // println!("{:?}", slice::FlexibleSlice::try_from(",").unwrap());

}