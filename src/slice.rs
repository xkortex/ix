use std::fmt;
use std::convert::TryFrom;
use std::error::Error;


extern crate regex;

#[derive(Debug, Clone, PartialEq, Eq)]
pub struct SliceError {
    kind: SliceErrorKind,
    msg: String,
}

#[derive(Debug, Clone, PartialEq, Eq)]
#[non_exhaustive]
enum SliceErrorKind {
    ParseSliceError,
    ParseIntError,
}

impl fmt::Display for SliceError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self.kind {
            SliceErrorKind::ParseSliceError => write!(f, "Unable to parse slice notation "),
            SliceErrorKind::ParseIntError => write!(f, "Unable to parse int in slice notation "),
        }
    }
}

impl SliceError {
    fn general_error() -> Self{
        SliceError{kind: SliceErrorKind::ParseSliceError, msg: String::from("General failure")}
    }
}

impl From<std::num::ParseIntError> for SliceError {
    fn from(e: std::num::ParseIntError) -> SliceError {
        SliceError{kind: SliceErrorKind::ParseIntError, msg: format!("{}", e.to_string())}
    }
}

#[derive(Debug, Clone)]
pub struct RawSlice {
    start: Option<String>,
    stop: Option<String>,
    step: Option<String>,
}

#[derive(Debug, Copy, Clone)]
pub struct FlexibleSlice {
    start: Option<isize>,
    stop: Option<isize>,
    step: Option<isize>,
}

#[derive(Debug, Clone)]
pub struct ProperSlicer {
    start: Option<usize>,
    stop: Option<usize>,
    step: usize,
    sep: regex::Regex
}

impl FlexibleSlice {
    fn new() -> Self {
        FlexibleSlice{ start: None, stop: None, step: None }
    }
}

impl From<String> for FlexibleSlice {
    fn from(s: String) -> Self {
        let parts: Vec<&str> = s.split(':').collect();
        println!(">{}< len: {}",s, parts.len());
        FlexibleSlice::new()
    }
}

// impl From<&str> for FlexibleSlice {
//     fn from(s: &str) -> Self {
//         let parts: Vec<&str> = s.split(':').collect();
//         println!(">{}< len: {}",s, parts.len());
//         FlexibleSlice::new()
//     }
// }

impl TryFrom<&str> for FlexibleSlice {
    type Error = SliceError;
    fn try_from(s: &str) -> Result<Self, SliceError> {
        // match s {
        //     "" => return Ok(FlexibleSlice::new()),
        //     ":" => return Ok(FlexibleSlice::new()),
        //     _ => {}
        // }
        let parts: Vec<&str> = s.split(':').collect();
        println!(">{}< len: {}",s, parts.len());
        match parts.len() {
            1 => {
                let start = match parts[0] {
                    "" => None,
                    x => Some(x.parse::<isize>()?),
                };
                let stop = start.map(|x| x + 1);
                Ok(FlexibleSlice{ start, stop, step: None })
            },
            2 => {
                let start = match parts[0] {
                    "" => None,
                    x => Some(x.parse::<isize>()?),
                };
                let stop = match parts[1] {
                    "" => None,
                    x => Some(x.parse::<isize>()?),
                };
                Ok(FlexibleSlice{ start, stop, step: None })
            },
            3 => {
                let start = match parts[0] {
                    "" => None,
                    x => Some(x.parse::<isize>()?),
                };
                let stop = match parts[1] {
                    "" => None,
                    x => Some(x.parse::<isize>()?),
                };
                let step = match parts[2] {
                    "" => None,
                    x => Some(x.parse::<isize>()?),
                };
                Ok(FlexibleSlice{ start, stop, step })
            },
            _ => Err(SliceError::general_error())
        }
    }
}