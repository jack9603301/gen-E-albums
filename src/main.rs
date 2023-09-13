extern crate mime;
use argparse::{ArgumentParser, Store};
use image::GenericImage;
use image::RgbaImage;
use std::fs;
use std::path::Path;
use std::process::Command;
use gostd::strings;
use image::GenericImageView;
use image::ImageBuffer;

fn image_process_mpv(buildroot: String, image: String, mut width: u32, mut height: u32) {
    println!("检测到图片：{}", image);
    println!("对图片文件 {} 的处理开始",image);

    // 获取图片基本信息
    let mut img = image::open(image.clone()).unwrap();
    let dimensions = img.dimensions();
    println!("原始图片大小： {:?}", dimensions);
    let raw_width = dimensions.0;
    let raw_height = dimensions.1;
    let mut save_width: u32 = width;
    let mut save_height: u32 = height;
    if (raw_width as f32 / raw_height as f32) > (width as f32 / height as f32) {
        save_height = save_width * raw_height / raw_width;
    } else {
        save_width = save_height * raw_width / raw_height;
    }

    //开始缩放处理
    img = img.resize(save_width, save_height, image::imageops::FilterType::Lanczos3);
    let save_dimensions = img.dimensions();
    println!("图片缩放后大小： {:?}", save_dimensions);
    save_width = save_dimensions.0;
    save_height = save_dimensions.1;

    // 计算偏移
    println!("计算图片缩放大小：{}x{}", save_width, save_height);
    let offsetx = ((width as i32 - save_width as i32) / 2) as u32;
    let offsety = ((height as i32 - save_height as i32) / 2) as u32;
    println!("计算贴图偏移量：({},{})", offsetx, offsety);
    

    // 开始图片转换

    let mut imgbuf: RgbaImage = ImageBuffer::from_pixel(width, height, image::Rgba([0, 0, 0, 0]));
    let _ = imgbuf.copy_from(&img, offsetx, offsety);

    let image_path = Path::new(&image);
    let image_filename = image_path.file_stem().unwrap().to_str().unwrap();

    let save_image_filename = buildroot + &String::from(image_filename) + &String::from(".png");

    let _ = imgbuf.save(save_image_filename.clone());
    println!("处理完成，图片文件保存至{}", save_image_filename);
}

fn main() {
    let mut root_path = "./".to_string();
    {
        let mut argparse = ArgumentParser::new();
        argparse.set_description("电子相册编译程序！");
        argparse.refer(&mut root_path)
            .add_argument("root", Store, "照片所在的目录");
        argparse.parse_args_or_exit();
    }

    let build_path = strings::TrimSuffix(root_path.as_str(), "/").to_owned() + "/build/";
    
    println!("检测到即将使用的build目录是：{}", build_path);
    if fs::metadata(build_path.clone()).is_err() {
        println!(">>>请注意，{} 目录不存在，现在创建！", build_path);
        let _ = fs::create_dir(build_path.clone());
    }
    
    
    println!("检测到您需要处理的文件是：{}", root_path);

    let mut images : Vec<String> = Vec::new();

    println!(">>>开始检查目录图片<<<");

    let paths = fs::read_dir(root_path).unwrap();
    for path in paths {
        let filename = String::from(path.unwrap().path().to_str().unwrap());
        let mime = match mime_guess::from_path(&filename).first() {
            Some(mime) => mime,
            None => mime::TEXT_PLAIN
        };
        println!("{} 开始检测：{}", filename, mime);
        if strings::HasPrefix(mime, "image/") {
            images.push(filename);
        }
    }
    println!(">>>检测完成<<<");

    for image in images {
        image_process_mpv(build_path.clone(), image, 1920, 1080);
    }
}
